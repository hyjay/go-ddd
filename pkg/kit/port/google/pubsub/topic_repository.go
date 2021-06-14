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

type TopicRepository struct {
	client *pubsub.Client
	topics *sync.Map
}

func NewTopicRepository(client *pubsub.Client) *TopicRepository {
	return &TopicRepository{
		client: client,
		topics: &sync.Map{},
	}
}

func (t *TopicRepository) GetOrCreate(ctx context.Context, topicID string) (*pubsub.Topic, error) {
	topic, ok := t.topics.Load(topicID)
	if !ok {
		var err error
		topic, err = t.createAndSaveTopic(ctx, topicID)
		if err != nil {
			return nil, err
		}
	}

	return topic.(*pubsub.Topic), nil
}

func (t *TopicRepository) createAndSaveTopic(ctx context.Context, topicID string) (*pubsub.Topic, error) {
	logrus.WithField("topic", topicID).Infof("Creating the topic")
	topic, err := t.client.CreateTopic(ctx, topicID)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() != codes.AlreadyExists {
			return nil, errors.Wrap(err, "failed to create the topic")
		}
		topic = t.client.Topic(topicID)
	}
	t.topics.Store(topicID, topic)
	return topic, nil
}
