// Code generated by mockery v0.0.0-dev.

package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SubscriptionRepositoryTestSuite struct {
	suite.Suite

	subscriptionRepository *SubscriptionRepository
	localPubSubServer      *LocalPubSubServer
	pubsubClient           *pubsub.Client
	fixedContext           context.Context
	fixedTopicID           string
	fixedSubscriptionID    string
	fixedTopic             *pubsub.Topic
}

func TestSubscriptionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionRepositoryTestSuite))
}

func (s *SubscriptionRepositoryTestSuite) TestGetOrCreate() {
	subscription, err := s.subscriptionRepository.GetOrCreate(s.fixedContext, s.fixedSubscriptionID, pubsub.SubscriptionConfig{Topic: s.fixedTopic})

	s.NoError(err)
	s.Equal(fmt.Sprint("projects/project/subscriptions/", s.fixedSubscriptionID), subscription.String())
}

func (s *SubscriptionRepositoryTestSuite) TestGetOrCreate_ShouldHandle_WhenSubscriptionAlreadyExist() {
	_, err := s.pubsubClient.CreateSubscription(s.fixedContext, s.fixedSubscriptionID, pubsub.SubscriptionConfig{Topic: s.fixedTopic})
	s.NoError(err)

	_, err = s.subscriptionRepository.GetOrCreate(s.fixedContext, s.fixedSubscriptionID, pubsub.SubscriptionConfig{Topic: s.fixedTopic})

	s.NoError(err)
}

func (s *SubscriptionRepositoryTestSuite) TestGetOrCreate_ShouldSaveCreatedSubscription() {
	subscription, err := s.subscriptionRepository.GetOrCreate(s.fixedContext, s.fixedSubscriptionID, pubsub.SubscriptionConfig{Topic: s.fixedTopic})
	s.NoError(err)

	savedSubscription, ok := s.subscriptionRepository.subscriptions.Load(s.fixedSubscriptionID)
	s.True(ok)
	s.Equal(subscription, savedSubscription)
}

func (s *SubscriptionRepositoryTestSuite) SetupTest() {
	s.localPubSubServer = NewLocalPubSubServer()
	var err error
	s.pubsubClient, err = s.localPubSubServer.CreateClient()
	s.NoError(err)

	s.subscriptionRepository = NewSubscriptionRepository(s.pubsubClient)
	s.fixedContext = context.Background()
	s.fixedTopicID = "FIXED_TOPIC"
	s.fixedSubscriptionID = "FIXED_SUBSCRIPTION"

	s.fixedTopic, err = s.pubsubClient.CreateTopic(s.fixedContext, s.fixedTopicID)
	s.NoError(err)
}

func (s *SubscriptionRepositoryTestSuite) TearDownTest() {
	err := s.localPubSubServer.TearDown()
	s.NoError(err)
}
