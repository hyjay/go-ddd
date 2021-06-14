package pubsub

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SubscriptionSchemeTestSuite struct {
	suite.Suite

	subscriptionScheme *SubscriptionScheme

	fixedTopicID     string
	fixedHandlerName string
}

func TestSubscriptionSchemeTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionSchemeTestSuite))
}

func (s *SubscriptionSchemeTestSuite) TestSubscriptionID() {
	subscriptionID := s.subscriptionScheme.SubscriptionID(s.fixedTopicID, s.fixedHandlerName)
	s.Equal("FIXED_TOPIC.subscriber.FIXED_HANDLER.v1", subscriptionID)
}

func (s *SubscriptionSchemeTestSuite) SetupTest() {
	s.subscriptionScheme = NewSubscriptionScheme("v1")
	s.fixedTopicID = "FIXED_TOPIC"
	s.fixedHandlerName = "FIXED_HANDLER"
}
