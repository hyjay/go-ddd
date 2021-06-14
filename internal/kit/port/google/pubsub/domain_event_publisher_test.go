package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DomainEventPublisherTestSuite struct {
	suite.Suite

	publisher *DomainEventPublisher

	localPubSubServer       *LocalPubSubServer
	messageChannel          chan *pubsub.Message
	fixedContext            context.Context
	fixedNamespace          string
	fixedPubsubSubscription *pubsub.Subscription
}

func TestDomainEventPublisherTestSuite(t *testing.T) {
	suite.Run(t, new(DomainEventPublisherTestSuite))
}

func (s *DomainEventPublisherTestSuite) TestPublish() {
	event := &greetedEvent{
		Greeting: "Hello world",
	}
	s.publisher.Publish(s.fixedContext, event)

	s.assertReceivedMessage(event)
}

func (s *DomainEventPublisherTestSuite) SetupTest() {
	s.fixedContext = context.Background()
	s.fixedNamespace = "FIXED_NAMESPACE"

	s.localPubSubServer = NewLocalPubSubServer()
	pubsubClient, err := s.localPubSubServer.CreateClient()
	s.NoError(err)

	topicRepository := NewTopicRepository(pubsubClient)
	topic, _ := topicRepository.GetOrCreate(s.fixedContext, fmt.Sprint(s.fixedNamespace, ".internal.greetedEvent.v1"))

	subscription, err := pubsubClient.CreateSubscription(s.fixedContext, "subscription", pubsub.SubscriptionConfig{Topic: topic})
	s.NoError(err)

	s.messageChannel = make(chan *pubsub.Message)
	go subscription.Receive(s.fixedContext, func(ctx context.Context, message *pubsub.Message) {
		s.messageChannel <- message
	})

	s.publisher = NewDomainEventPublisher(NewTopicScheme(s.fixedNamespace, "v1"), topicRepository)
}

func (s *DomainEventPublisherTestSuite) TearDownTest() {
	err := s.localPubSubServer.TearDown()
	s.NoError(err)
}

func (s *DomainEventPublisherTestSuite) assertReceivedMessage(event *greetedEvent) {
	msg := <-s.messageChannel
	expectedJSON, _ := json.Marshal(event)
	s.JSONEq(string(expectedJSON), string(msg.Data))
}
