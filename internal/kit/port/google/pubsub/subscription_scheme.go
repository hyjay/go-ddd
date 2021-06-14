package pubsub

import (
	"fmt"
)

type SubscriptionScheme struct {
	version string
}

func NewSubscriptionScheme(version string) *SubscriptionScheme {
	return &SubscriptionScheme{version: version}
}

func (s *SubscriptionScheme) SubscriptionID(topicID string, handlerName string) string {
	return fmt.Sprintf("%s.subscriber.%s.%s", topicID, handlerName, s.version)
}
