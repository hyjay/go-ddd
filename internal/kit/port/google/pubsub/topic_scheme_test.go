package pubsub

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TopicSchemeTestSuite struct {
	suite.Suite

	topicScheme *TopicScheme
}

func TestTopicSchemeTestSuite(t *testing.T) {
	suite.Run(t, new(TopicSchemeTestSuite))
}

func (s *TopicSchemeTestSuite) TestTopicID() {
	topicID := s.topicScheme.TopicID(&greetedEvent{})
	s.Equal("kit.internal.greetedEvent.v1", topicID)
}

func (s *TopicSchemeTestSuite) SetupTest() {
	s.topicScheme = NewTopicScheme("kit", "v1")
}
