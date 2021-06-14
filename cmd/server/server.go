package main

import (
	"github.com/emicklei/go-restful"
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
	userRepository := fakes.NewUserRepository()
	passwordHashService := bcrypt.NewPasswordHashService()

	container := restful.NewContainer()
	container.EnableContentEncoding(true)
	svc := service.NewService(userRepository, passwordHashService)
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
