package main

import (
	googpubsub "cloud.google.com/go/pubsub"
	"context"
	"github.com/emicklei/go-restful"
	"github.com/hyjay/go-ddd/internal/kit/port/google/pubsub"
	"github.com/hyjay/go-ddd/pkg/app/eventhandlers"
	"github.com/hyjay/go-ddd/pkg/app/service"
	"github.com/hyjay/go-ddd/pkg/domain/fakes"
	"github.com/hyjay/go-ddd/pkg/port/bcrypt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	httpServerAddress = "localhost:8080"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	localPubSubServer := pubsub.NewLocalPubSubServer()
	pubsubClient, err := localPubSubServer.CreateClient()
	if err != nil {
		logrus.WithError(err).Fatalf("error in creating a Pub/Sub client")
	}
	topicRepository := pubsub.NewTopicRepository(pubsubClient)
	topicScheme := pubsub.NewTopicScheme("account", "v1")
	subscriptionRepository := pubsub.NewSubscriptionRepository(pubsubClient)
	subscriptionScheme := pubsub.NewSubscriptionScheme("v1")
	domainEventPublisher := pubsub.NewDomainEventPublisher(topicScheme, topicRepository)
	domainEventSubscriber := pubsub.NewDomainEventSubscriber(topicScheme, subscriptionScheme, topicRepository, subscriptionRepository)

	userRepository := fakes.NewUserRepository()
	passwordHashService := bcrypt.NewPasswordHashService()

	sendWelcomeEmailWhenUserSignedUpEventHandler := eventhandlers.NewSendWelcomeEmailWhenUserSignedUpEventHandler()
	err = domainEventSubscriber.Subscribe(ctx, sendWelcomeEmailWhenUserSignedUpEventHandler, googpubsub.SubscriptionConfig{
		AckDeadline:         10 * time.Second,
		RetainAckedMessages: false,
		RetryPolicy: &googpubsub.RetryPolicy{
			MinimumBackoff: 10 * time.Second,
			MaximumBackoff: 60 * time.Second,
		},
	})
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"event_handler": sendWelcomeEmailWhenUserSignedUpEventHandler.Name()}).
			Fatalf("error in subscribing the domain event")
	}

	container := restful.NewContainer()
	container.EnableContentEncoding(true)
	svc := service.NewService(userRepository, passwordHashService, domainEventPublisher)
	svc.Register(container)

	runHTTPServer(container)
}

func runHTTPServer(container *restful.Container) {
	logrus.Infof("Starting HTTP server at %s", httpServerAddress)
	if err := http.ListenAndServe(httpServerAddress, container); err != nil {
		logrus.
			WithError(errors.Wrapf(err, "failed to listen HTTP server")).
			Fatalf("failed to listen HTTP server")
	}
}
