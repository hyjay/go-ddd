package service

import (
	"context"
	"github.com/emicklei/go-restful"
	"github.com/hyjay/go-ddd/pkg/domain"
	"net/http"
)

type PasswordHashService interface {
	HashPassword(ctx context.Context, plain string) (string, error)
}

type Service struct {
	userRepository      domain.UserRepository
	passwordHashService PasswordHashService
}

func NewService(userRepository domain.UserRepository, passwordHashService PasswordHashService) *Service {
	return &Service{userRepository: userRepository, passwordHashService: passwordHashService}
}

type user struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type signupRequest struct {
	User *user `json:"user"`
}

func (s *Service) Register(container *restful.Container) {
	service := &restful.WebService{}
	service.Path("/v1/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	service.Route(
		service.POST("").
			To(s.signup).
			Reads(signupRequest{}).
			Writes(user{}))

	container.Add(service)
}

func (s *Service) signup(request *restful.Request, response *restful.Response) {
	ctx := request.Request.Context()
	signupRequest := &signupRequest{}
	if err := request.ReadEntity(signupRequest); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}
	userResource := signupRequest.User
	hashedPassword, err := s.passwordHashService.HashPassword(ctx, userResource.Password)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	userEntity := makeUserEntity(signupRequest.User, "SOME_RANDOM_UUID", hashedPassword)
	if err := s.userRepository.Save(ctx, userEntity); err != nil {
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
