package pubsub

import (
	"fmt"
	"github.com/hyjay/go-ddd/internal/kit"
)

type TopicScheme struct {
	namespace string
	version   string
}

func NewTopicScheme(namespace string, version string) *TopicScheme {
	return &TopicScheme{namespace: namespace, version: version}
}

func (t *TopicScheme) TopicID(event kit.DomainEvent) string {
	return fmt.Sprintf("%s.internal.%s.%s", t.namespace, event.Name(), t.version)
}
