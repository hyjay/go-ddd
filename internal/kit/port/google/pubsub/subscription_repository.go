package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type SubscriptionRepository struct {
	client        *pubsub.Client
	subscriptions *sync.Map
}

func NewSubscriptionRepository(client *pubsub.Client) *SubscriptionRepository {
	return &SubscriptionRepository{
		client:        client,
		subscriptions: &sync.Map{},
	}
}

func (sr *SubscriptionRepository) GetOrCreate(ctx context.Context, subscriptionID string, config pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	subscription, ok := sr.subscriptions.Load(subscriptionID)
	if !ok {
		var err error
		subscription, err = sr.createAndSaveSubscription(ctx, subscriptionID, config)
		if err != nil {
			return nil, err
		}
	}
	return subscription.(*pubsub.Subscription), nil
}

func (sr *SubscriptionRepository) createAndSaveSubscription(
	ctx context.Context, subscriptionID string, config pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	logrus.WithField("subscription", subscriptionID).Infof("Creating the subscription")
	subscription, err := sr.client.CreateSubscription(ctx, subscriptionID, config)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() != codes.AlreadyExists {
			return nil, errors.Wrapf(err, "failed to create the subscription")
		}
		subscription = sr.client.Subscription(subscriptionID)
	}

	sr.subscriptions.Store(subscriptionID, subscription)
	return subscription, nil
}
