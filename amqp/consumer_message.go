package amqp

import (
	"github.com/marcuzy/rabbus"
	"time"
)

func newConsumerMessage(m *rabbus.ConsumerMessage) ConsumerMessage {
	return &consumerMessage{
		rabbusConsumerMessage: m,

		ContentType:     m.ContentType,
		ContentEncoding: m.ContentEncoding,
		DeliveryMode:    m.DeliveryMode,
		Priority:        m.Priority,
		CorrelationId:   m.CorrelationId,
		ReplyTo:         m.ReplyTo,
		Expiration:      m.Expiration,
		MessageId:       m.MessageId,
		Timestamp:       m.Timestamp,
		Type:            m.Type,
		ConsumerTag:     m.ConsumerTag,
		MessageCount:    m.MessageCount,
		DeliveryTag:     m.DeliveryTag,
		Redelivered:     m.Redelivered,
		Exchange:        m.Exchange,
		Headers:         m.Headers,
		Key:             m.Key,
		Body:            m.Body,
	}
}

type (
	ConsumerMessage interface {
		Ack(multiple bool) error
		Nack(multiple, requeue bool) error
		Reject(requeue bool) error

		GetContentType() string
		GetContentEncoding() string
		// DeliveryMode queue implementation use, non-persistent (1) or persistent (2)
		GetDeliveryMode() uint8
		// Priority queue implementation use, 0 to 9
		GetPriority() uint8
		// CorrelationId application use, correlation identifier
		GetCorrelationId() string
		// ReplyTo application use, address to to reply to (ex: RPC)
		GetReplyTo() string
		// Expiration implementation use, message expiration spec
		GetExpiration() string
		// MessageId application use, message identifier
		GetMessageId() string
		// Timestamp application use, message timestamp
		GetTimestamp() time.Time
		// Type application use, message type name
		GetType() string
		// ConsumerTag valid only with Channel.Consume
		GetConsumerTag() string
		// MessageCount valid only with Channel.Get
		GetMessageCount() uint32
		GetDeliveryTag() uint64
		GetRedelivered() bool
		GetExchange() string
		// Headers application or header exchange table
		GetHeaders() map[string]interface{}
		// Key basic.publish routing key
		GetKey() string
		GetBody() []byte
	}

	// consumerMessage is a wrapper around lib's message to decouple things
	consumerMessage struct {
		rabbusConsumerMessage *rabbus.ConsumerMessage

		ContentType     string
		ContentEncoding string
		DeliveryMode    uint8
		Priority        uint8
		CorrelationId   string
		ReplyTo         string
		Expiration      string
		MessageId       string
		Timestamp       time.Time
		Type            string
		ConsumerTag     string
		MessageCount    uint32
		DeliveryTag     uint64
		Redelivered     bool
		Exchange        string
		Headers         map[string]interface{}
		Key             string
		Body            []byte
	}
)

func (cm *consumerMessage) GetContentType() string {
	return cm.ContentType
}

func (cm *consumerMessage) GetContentEncoding() string {
	return cm.ContentEncoding
}

func (cm *consumerMessage) GetDeliveryMode() uint8 {
	return cm.DeliveryMode
}

func (cm *consumerMessage) GetPriority() uint8 {
	return cm.Priority
}

func (cm *consumerMessage) GetCorrelationId() string {
	return cm.CorrelationId
}

func (cm *consumerMessage) GetReplyTo() string {
	return cm.ReplyTo
}

func (cm *consumerMessage) GetExpiration() string {
	return cm.Expiration
}

func (cm *consumerMessage) GetMessageId() string {
	return cm.MessageId
}

func (cm *consumerMessage) GetTimestamp() time.Time {
	return cm.Timestamp
}

func (cm *consumerMessage) GetType() string {
	return cm.Type
}

func (cm *consumerMessage) GetConsumerTag() string {
	return cm.ConsumerTag
}

func (cm *consumerMessage) GetMessageCount() uint32 {
	return cm.MessageCount
}

func (cm *consumerMessage) GetDeliveryTag() uint64 {
	return cm.DeliveryTag
}

func (cm *consumerMessage) GetRedelivered() bool {
	return cm.Redelivered
}

func (cm *consumerMessage) GetExchange() string {
	return cm.Exchange
}

func (cm *consumerMessage) GetHeaders() map[string]interface{} {
	return cm.Headers
}

func (cm *consumerMessage) GetKey() string {
	return cm.Key
}

func (cm *consumerMessage) GetBody() []byte {
	return cm.Body
}

func (cm *consumerMessage) Ack(multiple bool) error {
	return cm.rabbusConsumerMessage.Ack(multiple)
}

func (cm *consumerMessage) Nack(multiple, requeue bool) error {
	return cm.rabbusConsumerMessage.Nack(multiple, requeue)
}

func (cm *consumerMessage) Reject(requeue bool) error {
	return cm.rabbusConsumerMessage.Reject(requeue)
}
