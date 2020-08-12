package amqp

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"bitbucket.org/iwlab-standuply/slackteams-api/shared"
	"github.com/marcuzy/rabbus"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNoRPCResponse = errors.New("no rpc response")
)

// NOTE: instance of rabbus.Rabbus uses single channel for all operations so can Produce only one message at once (blocking until prev message is sent)
// maybe we should create an AMQP instance (lib's internal class) and pass it to every Robbus instance to share connection
type (
	// Client provides easy robust communication with RabbitMQ.
	// Here we use our own types and abstractions to decouple with any underlying lib.
	// It will let us replace libs without affecting most of the codebase.
	Client interface {
		Consume(params ConsumeParams) (<-chan ConsumerMessage, error)
		Produce(msg Message) error
		Publish(ctx context.Context, params ResponseParams) error
		Request(ctx context.Context, params RequestParams) ([]byte, error)

		ConsumeRPCRequests(routingKey string) (<-chan ConsumerMessage, error)
		PublishRPCResponse(ctx context.Context, params RPCResponseParams) error

		Connect(ctx context.Context) error
		Close()
	}

	ConsumeParams struct {
		Exchange   string
		Kind       string
		RoutingKey string
		Queue      string
	}

	RequestParams struct {
		Exchange   string
		RoutingKey string
		Payload    interface{}
	}

	ResponseParams struct {
		Exchange   string
		RoutingKey string
		MessageID  string
		Payload    interface{}
	}

	RPCResponseParams struct {
		RoutingKey string
		MessageID  string
		Payload    interface{}
	}

	amqpClient struct {
		url                string
		rpcResponsesQ      string
		rpcResponseTimeout time.Duration

		r            *rabbus.Rabbus
		rpcResponses map[string]chan ConsumerMessage
	}
)

func NewClient(url string) Client {
	r, err := rabbus.New(
		url,
		rabbus.Durable(true),
		rabbus.Attempts(5),
		rabbus.Sleep(time.Second*2),
		rabbus.Threshold(3),
		// rabbus.OnStateChange(cbStateChangeFunc),
	)

	if err != nil {
		panic(err)
	}

	return &amqpClient{
		url:                url,
		rpcResponses:       make(map[string]chan ConsumerMessage),
		rpcResponsesQ:      "slackTeams.rpcResponses",
		rpcResponseTimeout: time.Minute * 1,
		r:                  r,
	}
}

func (c *amqpClient) Connect(ctx context.Context) error {
	go func() {
		err := c.r.Run(context.Background())

		if err != nil {
			log.WithError(err).Error()
		}
	}()

	err := c.startListeningToRPCResponses()

	if err != nil {
		return err
	}

	return nil
}

func (c *amqpClient) Close() {
	if c.r != nil {
		c.r.Close()
	}
}

func (c *amqpClient) Publish(ctx context.Context, params ResponseParams) error {
	// messageId := shared.RandStringBytesMaskImprSrcUnsafe(24)

	payload, err := json.Marshal(params.Payload)

	if err != nil {
		return err
	}

	err = c.Produce(&message{
		Exchange:    params.Exchange,
		Key:         params.RoutingKey,
		MessageId:   params.MessageID,
		Kind:        "direct",
		ContentType: "application/json",
		Payload:     payload,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *amqpClient) PublishRPCResponse(ctx context.Context, params RPCResponseParams) error {
	if params.RoutingKey == "" || params.MessageID == "" {
		return errors.New("Empty params")
	}

	return c.Publish(ctx, ResponseParams{
		Exchange:   "slackTeams.api.response",
		RoutingKey: params.RoutingKey,
		MessageID:  params.MessageID,
		Payload:    params.Payload,
	})
}

func (c *amqpClient) Request(ctx context.Context, params RequestParams) ([]byte, error) {
	messageId := shared.RandStringBytesMaskImprSrcUnsafe(24)

	payload, err := json.Marshal(params.Payload)

	if err != nil {
		return nil, err
	}

	err = c.Produce(&message{
		Exchange:    params.Exchange,
		Key:         params.RoutingKey,
		MessageId:   messageId,
		Kind:        "topic",
		ContentType: "application/json",
		Payload:     payload,
		ReplyTo:     c.rpcResponsesQ,
	})

	if err != nil {
		return nil, err
	}

	responseCtx, cancel := context.WithTimeout(ctx, c.rpcResponseTimeout)
	defer cancel()

	m, err := c.waitRPCResponse(responseCtx, messageId)

	if err != nil {
		return nil, err
	}

	return m.GetBody(), nil
}

// TODO: add ctx to cancel for graceful shutdown
func (c *amqpClient) Consume(params ConsumeParams) (<-chan ConsumerMessage, error) {
	msgs, err := c.r.Listen(rabbus.ListenConfig{
		Exchange: params.Exchange,
		Key:      params.RoutingKey,
		Kind:     params.Kind,
		Queue:    params.Queue,
	})

	if err != nil {
		return nil, err
	}

	res := make(chan ConsumerMessage, 256)

	go func() {
		for {
			m, ok := <-msgs
			if !ok {
				close(res)
				break
			}
			res <- newConsumerMessage(&m)
		}
	}()

	return res, nil
}

func (c *amqpClient) ConsumeRPCRequests(routingKey string) (<-chan ConsumerMessage, error) {
	return c.Consume(ConsumeParams{
		Exchange:   "slackTeams.api.tx",
		RoutingKey: routingKey,
		Queue:      "slackTeams.api." + routingKey,
		Kind:       "topic",
	})
}

// TODO fix
func (c *amqpClient) Produce(msg Message) error {
	c.r.EmitAsync() <- newRobbusMessage(msg)

	select {
	case err := <-c.r.EmitErr():
		return err
	case <-c.r.EmitOk():
		return nil
	}
}

func (c *amqpClient) waitRPCResponse(ctx context.Context, messageID string) (ConsumerMessage, error) {
	chRes := make(chan ConsumerMessage)
	defer close(chRes)

	// Register a channel to receive a response to the request
	c.rpcResponses[messageID] = chRes

	select {
	case res := <-chRes:
		return res, nil
	case <-ctx.Done():
		delete(c.rpcResponses, messageID)

		return nil, ctx.Err()
	}
}

func (c *amqpClient) startListeningToRPCResponses() error {
	msgs, err := c.Consume(ConsumeParams{
		Exchange:   "experts.api.response",
		Kind:       "direct",
		RoutingKey: c.rpcResponsesQ,
		Queue:      c.rpcResponsesQ,
	})

	if err != nil {
		return err
	}

	go func() {
		for m := range msgs {
			ch, has := c.rpcResponses[m.GetCorrelationId()]

			if !has {
				// TODO log about a message will be lost
				if err := m.Ack(false); err != nil {
					// TODO what to do?
				}

				continue
			}

			ch <- m

			if err := m.Ack(false); err != nil {
				log.WithError(err).WithField("message", m).Error("Failed to ack rpc response")
			}

			delete(c.rpcResponses, m.GetCorrelationId())
		}

		for corrId, ch := range c.rpcResponses {
			close(ch)

			delete(c.rpcResponses, corrId)
		}
	}()

	return nil
}

func newRobbusMessage(m Message) rabbus.Message {
	return rabbus.Message{
		Exchange:        m.GetExchange(),
		Kind:            m.GetKind(),
		Key:             m.GetKey(),
		Payload:         m.GetPayload(),
		DeliveryMode:    m.GetDeliveryMode(),
		ContentType:     m.GetContentType(),
		Headers:         m.GetHeaders(),
		ReplyTo:         m.GetReplyTo(),
		ContentEncoding: m.GetContentEncoding(),
		MessageId:       m.GetMessageId(),
	}
}
