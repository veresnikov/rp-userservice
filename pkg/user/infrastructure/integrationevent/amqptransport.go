package integrationevent

import (
	"context"
	"encoding/json"
	"errors"

	"gitea.xscloud.ru/xscloud/golib/pkg/application/logging"
	"gitea.xscloud.ru/xscloud/golib/pkg/infrastructure/amqp"
)

func NewAMQPTransport(logger logging.Logger) AMQPTransport {
	return &amqpTransport{
		logger: logger,
	}
}

type AMQPTransport interface {
	Handler() amqp.Handler
}

type amqpTransport struct {
	logger logging.Logger
}

func (t *amqpTransport) Handler() amqp.Handler {
	return t.handle
}

func (t *amqpTransport) handle(ctx context.Context, delivery amqp.Delivery) error {
	l := t.logger.WithFields(logging.Fields{
		"routing_key":    delivery.RoutingKey,
		"correlation_id": delivery.CorrelationID,
		"content_type":   delivery.ContentType,
	})
	if delivery.ContentType != ContentType {
		l.Warning(errors.New("invalid content type"), "skipping")
		return nil
	}
	l = l.WithField("body", json.RawMessage(delivery.Body))

	var err error

	if err != nil {
		l.Error(err, "failed to handle message")
	} else {
		l.Info("successfully handled message")
	}
	return err
}
