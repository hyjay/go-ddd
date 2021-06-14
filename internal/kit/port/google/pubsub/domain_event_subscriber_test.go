package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/hyjay/go-ddd/internal/kit"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DomainEventHandlerRegistryTestSuite struct {
	suite.Suite

	subscriber *DomainEventSubscriber

	publisher         *DomainEventPublisher
	handlerSpy        *handlerSpy
	localPubSubServer *LocalPubSubServer
	fixedContext      context.Context
}

type handlerSpy struct {
	receivedEventChannel chan kit.DomainEvent
}

func newHandlerSpy() *handlerSpy {
	return &handlerSpy{receivedEventChannel: make(chan kit.DomainEvent, 1)}
}

func (h *handlerSpy) TargetEvent() kit.DomainEvent {
	return &greetedEvent{}
}

func (h *handlerSpy) Name() string {
	return "handlerSpy"
}

func (h *handlerSpy) Handle(ctx context.Context, event kit.DomainEvent) error {
	h.receivedEventChannel <- event
	return nil
}

func (h *handlerSpy) Poll() kit.DomainEvent {
	return <-h.receivedEventChannel
}

func TestDomainEventHandlerRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(DomainEventHandlerRegistryTestSuite))
}

func (s *DomainEventHandlerRegistryTestSuite) TestSubscribe() {
	event := &greetedEvent{"Hello world"}
	s.publisher.Publish(context.Background(), event)

	handledEvent := s.handlerSpy.Poll()
	s.Equal(event, handledEvent)
}

func (s *DomainEventHandlerRegistryTestSuite) SetupTest() {
	s.localPubSubServer = NewLocalPubSubServer()
	pubsubClient, err := s.localPubSubServer.CreateClient()
	s.NoError(err)
	topicScheme := NewTopicScheme("FIXED_NAMESPACE", "v1")
	subscriptionScheme := NewSubscriptionScheme("v1")
	topicRepository := NewTopicRepository(pubsubClient)
	s.publisher = NewDomainEventPublisher(topicScheme, topicRepository)
	s.subscriber = NewDomainEventSubscriber(topicScheme, subscriptionScheme, topicRepository, NewSubscriptionRepository(pubsubClient))
	s.handlerSpy = newHandlerSpy()
	s.fixedContext = context.Background()

	err = s.subscriber.Subscribe(s.fixedContext, s.handlerSpy, pubsub.SubscriptionConfig{})
	s.NoError(err)

}

func (s *DomainEventHandlerRegistryTestSuite) TearDownTest() {
	err := s.localPubSubServer.TearDown()
	s.NoError(err)
}
