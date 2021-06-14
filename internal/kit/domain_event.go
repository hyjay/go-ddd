package kit

import (
	"context"
)

type DomainEvent interface {
	Name() string
}

type DomainEventHandler interface {
	TargetEvent() DomainEvent
	Name() string
	Handle(ctx context.Context, event DomainEvent) error
}

type DomainEventPublisher interface {
	Publish(ctx context.Context, event DomainEvent)
}
