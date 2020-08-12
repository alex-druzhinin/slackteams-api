package amqp

type (
	Message interface {
		// Exchange the exchange name.
		GetExchange() string
		// Kind the exchange type.
		GetKind() string
		// Key the routing key name.
		GetKey() string
		// Payload the message payload.
		GetPayload() []byte
		// DeliveryMode indicates if the is Persistent or Transient.
		GetDeliveryMode() uint8
		// ContentType the message content-type.
		GetContentType() string
		// Headers the message application headers
		GetHeaders() map[string]interface{}
		// ContentEncoding the message encoding.
		GetContentEncoding() string
		GetMessageId() string
		GetReplyTo() string
	}

	message struct {
		Exchange        string
		Kind            string
		Key             string
		Payload         []byte
		DeliveryMode    uint8
		ContentType     string
		Headers         map[string]interface{}
		ContentEncoding string
		MessageId       string
		ReplyTo         string
	}
)

func (m message) GetExchange() string {
	return m.Exchange
}

func (m message) GetKind() string {
	return m.Kind
}

func (m message) GetKey() string {
	return m.Key
}

func (m message) GetPayload() []byte {
	return m.Payload
}

func (m message) GetDeliveryMode() uint8 {
	return m.DeliveryMode
}

func (m message) GetContentType() string {
	return m.ContentType
}

func (m message) GetHeaders() map[string]interface{} {
	return m.Headers
}

func (m message) GetContentEncoding() string {
	return m.ContentEncoding
}

func (m message) GetMessageId() string {
	return m.MessageId
}

func (m message) GetReplyTo() string {
	return m.ReplyTo
}
