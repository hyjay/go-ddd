package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/hyjay/go-ddd/internal/kit"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DomainEventSubscriber struct {
	topicScheme            *TopicScheme
	subscriptionScheme     *SubscriptionScheme
	topicRepository        *TopicRepository
	subscriptionRepository *SubscriptionRepository
}

func NewDomainEventSubscriber(topicScheme *TopicScheme, subscriptionScheme *SubscriptionScheme, topicRepository *TopicRepository, subscriptionRepository *SubscriptionRepository) *DomainEventSubscriber {
	return &DomainEventSubscriber{topicScheme: topicScheme, subscriptionScheme: subscriptionScheme, topicRepository: topicRepository, subscriptionRepository: subscriptionRepository}
}

func (d *DomainEventSubscriber) Subscribe(ctx context.Context, handler kit.DomainEventHandler, config pubsub.SubscriptionConfig) error {
	logger := logrus.WithFields(logrus.Fields{
		"event_handler": handler.Name(),
		"event":         handler.TargetEvent().Name(),
	})
	topic, err := d.topicRepository.GetOrCreate(ctx, d.topicID(handler))
	if err != nil {
		return err
	}

	config.Topic = topic
	subscriptionID := d.subscriptionID(handler)
	subscription, err := d.subscriptionRepository.GetOrCreate(ctx, subscriptionID, config)
	if err != nil {
		return err
	}

	go func() {
		ctx := context.Background()
		err := subscription.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
			event := handler.TargetEvent()
			if err := d.unmarshalData(event, message.Data); err != nil {
				logger.WithError(err).
					WithFields(logrus.Fields{"payload": string(message.Data)}).
					Errorf("Failed to unmarshal the Pub/Sub message")
				message.Nack()
				return
			}

			if err = handler.Handle(ctx, event); err != nil {
				logger.WithError(err).
					WithFields(logrus.Fields{"payload": string(message.Data)}).
					Errorf("Failed to handle the domain event")
				message.Nack()
				return
			}

			message.Ack()
		})
		if err != nil {
			// err returned by Receive() is a non-retryable error
			logger.WithError(err).WithFields(logrus.Fields{"subscription": subscriptionID}).
				Fatalf("The Pub/Sub subscriber got a non-retryable error")
		}
	}()

	return nil
}

func (d *DomainEventSubscriber) subscriptionID(handler kit.DomainEventHandler) string {
	return d.subscriptionScheme.SubscriptionID(d.topicID(handler), handler.Name())
}

func (d *DomainEventSubscriber) topicID(handler kit.DomainEventHandler) string {
	return d.topicScheme.TopicID(handler.TargetEvent())
}

func (d *DomainEventSubscriber) unmarshalData(event kit.DomainEvent, data []byte) error {
	if err := json.Unmarshal(data, event); err != nil {
		return errors.Wrapf(err, "error in unmarshalling")
	}

	return nil
}
