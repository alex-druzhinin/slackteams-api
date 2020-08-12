package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"bitbucket.org/iwlab-standuply/slackteams-api/amqp"
	log "github.com/sirupsen/logrus"
)

var (
	ErrStanduplyClientResponseType = errors.New("standuply response has invalid type")
)

type Services struct {
}

type Server interface {
	Run() error
}

func NewTeamsRPCServer(amqpClient amqp.Client, repo SlackTeamsRepository) Server {
	return &rpcServer{
		c:    amqpClient,
		repo: repo,
	}
}

type rpcServer struct {
	c    amqp.Client
	repo SlackTeamsRepository
}

type getTeamByIDRequest struct {
	TeamID string `json:"teamId"`
}

func (s *rpcServer) Run() error {
	if err := s.observeGetTeam(); err != nil {
		return err
	}

	return nil
}

type simpleResponse struct {
	OK    bool    `json:"ok"`
	Error *string `json:"error,omitempty"`
}

type teamResponse struct {
	OK    bool      `json:"ok"`
	Error *string   `json:"error,omitempty"`
	Data  *SlackTeam `json:"data"`
}

func (s *rpcServer) observeGetTeam() error {
	messages, err := s.c.ConsumeRPCRequests("getTeam")

	if err != nil {
		return err
	}

	go func() {
		for m := range messages {
			s.handleSafely(s.handleGetTeam, m)

			err := m.Ack(false)
			if err != nil {
				log.WithError(err).Errorf("Failed to ack message %+v", m)
			}
			log.Debugf("Message ack")
		}
	}()

	return nil
}

func (s *rpcServer) handleGetTeam(m amqp.ConsumerMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	payload := teamResponse{
		OK: true,
	}

	var c getTeamByIDRequest

	err := json.Unmarshal(m.GetBody(), &c)

	if err != nil {
		s.responseWithError(ctx, m, err, "Failed to unmarshall getTeam event")
		return
	}

	t, err := s.repo.FindTeamByID(ctx, c.TeamID)
	if err != nil {
		s.responseWithError(ctx, m, err, "Failed to Find team")
		return
	}
	log.Debugf("FindTeamByID %s - %v", c.TeamID, t)

	payload.Data = t
	s.response(ctx, m, payload)
}

func (s *rpcServer) response(ctx context.Context, message amqp.ConsumerMessage, payload interface{}) {
	err := s.c.PublishRPCResponse(ctx, amqp.RPCResponseParams{
		RoutingKey: message.GetReplyTo(),
		MessageID:  message.GetCorrelationId(),
		Payload:    payload,
	})
	if err != nil {
		log.WithError(err).Errorf("Failed to publish response for message %+v with payload %+v", message, payload)
	}
}

func (s *rpcServer) responseWithError(ctx context.Context, message amqp.ConsumerMessage, err error, msg string) {
	log.WithError(err).Error(msg)

	errMessage := err.Error()
	payload := simpleResponse{
		OK:    false,
		Error: &errMessage,
	}

	s.response(ctx, message, payload)
}

func (s *rpcServer) handleSafely(f func(m amqp.ConsumerMessage), msg amqp.ConsumerMessage) {
	defer func() {
		if r := recover(); r != nil {
			var err error

			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}

			log.WithError(err).WithField("amqpMsg", msg).Error("panic in RpcServer recovered")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
			defer cancel()

			s.responseWithError(ctx, msg, err, "fatal error")
		}
	}()

	f(msg)
}
