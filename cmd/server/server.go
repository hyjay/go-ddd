package main

import (
	"github.com/emicklei/go-restful"
	"github.com/hyjay/go-ddd/internal/kit/port/google/pubsub"
	"github.com/hyjay/go-ddd/pkg/app/service"
	"github.com/hyjay/go-ddd/pkg/domain/fakes"
	"github.com/hyjay/go-ddd/pkg/port/bcrypt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	httpServerAddress = "localhost:8080"
)

func main() {
	localPubSubServer := pubsub.NewLocalPubSubServer()
	pubsubClient, err := localPubSubServer.CreateClient()
	if err != nil {
		logrus.WithError(err).Fatalf("error in creating a Pub/Sub client")
	}
	topicRepository := pubsub.NewTopicRepository(pubsubClient)
	topicScheme := pubsub.NewTopicScheme("account", "v1")
	domainEventPublisher := pubsub.NewDomainEventPublisher(topicScheme, topicRepository)

	userRepository := fakes.NewUserRepository()
	passwordHashService := bcrypt.NewPasswordHashService()

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
