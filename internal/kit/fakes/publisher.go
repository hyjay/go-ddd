package fakes

import (
	"context"
	"github.com/hyjay/go-ddd/internal/kit"
)

type DomainEventPublisher struct {
	publishedEvents []kit.DomainEvent
}

func NewDomainEventPublisher() *DomainEventPublisher {
	return &DomainEventPublisher{}
}

func (p *DomainEventPublisher) Publish(ctx context.Context, event kit.DomainEvent) {
	p.publishedEvents = append(p.publishedEvents, event)
}

func (p *DomainEventPublisher) PublishedEvents() []kit.DomainEvent {
	return p.publishedEvents
}
