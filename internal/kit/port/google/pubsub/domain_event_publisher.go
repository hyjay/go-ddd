package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/hyjay/go-ddd/internal/kit"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DomainEventPublisher struct {
	topicScheme     *TopicScheme
	topicRepository *TopicRepository
}

func NewDomainEventPublisher(topicScheme *TopicScheme, topicRepository *TopicRepository) *DomainEventPublisher {
	return &DomainEventPublisher{topicScheme: topicScheme, topicRepository: topicRepository}
}

func (p *DomainEventPublisher) Publish(ctx context.Context, evt kit.DomainEvent) {
	if err := p.publish(ctx, evt); err != nil {
		logrus.WithField("event", evt.Name()).WithField("payload", evt).
			WithError(err).Errorf("Failed to publish the Pub/Sub message")
		return
	}
	return
}

func (p *DomainEventPublisher) publish(ctx context.Context, evt kit.DomainEvent) error {
	topicID := p.topicScheme.TopicID(evt)
	t, err := p.topicRepository.GetOrCreate(ctx, topicID)
	if err != nil {
		return err
	}
	encoded, err := p.encodeMessage(evt)
	if err != nil {
		return err
	}
	t.Publish(ctx, &pubsub.Message{Data: encoded})
	return nil
}

func (p *DomainEventPublisher) encodeMessage(msg interface{}) ([]byte, error) {
	res, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode the domain event")
	}

	return res, nil
}
