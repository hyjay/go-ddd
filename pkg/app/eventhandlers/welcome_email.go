package eventhandlers

import (
	"context"
	"github.com/hyjay/go-ddd/internal/kit"
	"github.com/hyjay/go-ddd/pkg/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SendWelcomeEmailWhenUserSignedUpEventHandler struct{}

func NewSendWelcomeEmailWhenUserSignedUpEventHandler() *SendWelcomeEmailWhenUserSignedUpEventHandler {
	return &SendWelcomeEmailWhenUserSignedUpEventHandler{}
}

func (h *SendWelcomeEmailWhenUserSignedUpEventHandler) TargetEvent() kit.DomainEvent {
	return &domain.UserSignedUpEvent{}
}

func (h *SendWelcomeEmailWhenUserSignedUpEventHandler) Name() string {
	return "SendWelcomeEmailWhenUserSignedUpEventHandler"
}

func (h *SendWelcomeEmailWhenUserSignedUpEventHandler) Handle(ctx context.Context, event kit.DomainEvent) error {
	userSignedUpEvent, ok := event.(*domain.UserSignedUpEvent)
	if !ok {
		return errors.New("the event is not UserSignedUpEvent")
	}
	return h.doHandle(userSignedUpEvent)
}

func (h *SendWelcomeEmailWhenUserSignedUpEventHandler) doHandle(event *domain.UserSignedUpEvent) error {
	// TODO: Send a welcome email to the user's email
	logrus.WithFields(logrus.Fields{"user_id": event.UserID, "email": event.Email}).Infof("Sending a welcome email")
	return nil
}
