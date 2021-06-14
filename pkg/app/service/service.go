package service

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hyjay/go-ddd/internal/kit"
	"github.com/hyjay/go-ddd/pkg/domain"
	"net/http"
	"time"
)

const (
	userIDParameter = "user_id"
)

type PasswordHashService interface {
	HashPassword(ctx context.Context, plain string) (string, error)
}

type Service struct {
	userRepository       domain.UserRepository
	passwordHashService  PasswordHashService
	domainEventPublisher kit.DomainEventPublisher
}

func NewService(userRepository domain.UserRepository, passwordHashService PasswordHashService, domainEventPublisher kit.DomainEventPublisher) *Service {
	return &Service{userRepository: userRepository, passwordHashService: passwordHashService, domainEventPublisher: domainEventPublisher}
}

type user struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *Service) Register(container *restful.Container) {
	service := &restful.WebService{}
	service.Path("/v1/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	service.Route(
		service.POST("").
			To(s.signup).
			Reads(user{}).
			Writes(user{}))
	service.Route(
		service.GET(fmt.Sprintf("{%s}", userIDParameter)).
			To(s.getUser).
			Writes(user{}))

	container.Add(service)
}

func (s *Service) signup(request *restful.Request, response *restful.Response) {
	ctx := request.Request.Context()
	user := &user{}
	if err := request.ReadEntity(user); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}
	hashedPassword, err := s.passwordHashService.HashPassword(ctx, user.Password)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	userID := "SOME_RANDOM_UUID"
	userEntity := makeUserEntity(user, domain.UserID(userID), hashedPassword)
	if err := s.userRepository.Save(ctx, userEntity); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	s.domainEventPublisher.Publish(ctx, &domain.UserSignedUpEvent{
		UserID:    domain.UserID(userID),
		Email:     userEntity.Email(),
		FirstName: userEntity.FirstName(),
		LastName:  userEntity.LastName(),
		Timestamp: time.Now(),
	})
	if err := response.WriteEntity(makeUserResource(userEntity)); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func (s *Service) getUser(request *restful.Request, response *restful.Response) {
	ctx := request.Request.Context()
	userID := request.PathParameter(userIDParameter)
	userEntity, err := s.userRepository.GetByID(ctx, domain.UserID(userID))
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	if err := response.WriteEntity(makeUserResource(userEntity)); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func makeUserEntity(u *user, id domain.UserID, hashedPassword string) *domain.User {
	entity := &domain.User{}
	entity.Write(id, u.Email, hashedPassword, u.FirstName, u.LastName)
	return entity
}

func makeUserResource(entity *domain.User) *user {
	user := &user{}
	entity.Read(func(id domain.UserID, email string, hashedPassword string, firstName string, lastName string) {
		user.ID = string(id)
		user.Email = email
		user.FirstName = firstName
		user.LastName = lastName
	})
	return user
}
