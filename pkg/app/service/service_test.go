package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/hyjay/go-ddd/pkg/app/service/mocks"
	"github.com/hyjay/go-ddd/pkg/domain"
	domainmocks "github.com/hyjay/go-ddd/pkg/domain/mocks"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite

	service                 *Service
	mockUserRepository      *domainmocks.UserRepository
	mockPasswordHashService *mocks.PasswordHashService
	container               *restful.Container
	responseRecorder        *httptest.ResponseRecorder

	fixedUserID         string
	fixedEmail          string
	fixedPassword       string
	fixedFirstName      string
	fixedLastName       string
	fixedHashedPassword string
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestSignup() {
	s.mockPasswordHashService.On("HashPassword", mock.Anything, s.fixedPassword).
		Return(s.fixedHashedPassword, nil)
	s.mockUserRepository.On("Save", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		user.Read(func(id domain.UserID, email string, hashedPassword string, firstName string, lastName string) {
			s.Equal(domain.UserID(s.fixedUserID), id)
			s.Equal(s.fixedEmail, email)
			s.Equal(s.fixedHashedPassword, hashedPassword)
			s.Equal(s.fixedFirstName, firstName)
			s.Equal(s.fixedLastName, lastName)
		})
		return true
	})).
		Return(nil)

	request := &user{
		Email:     s.fixedEmail,
		Password:  s.fixedPassword,
		FirstName: s.fixedFirstName,
		LastName:  s.fixedLastName,
	}
	s.serve(http.MethodPost, "/v1/users", request)

	s.Equal(200, s.responseRecorder.Code)
	s.equalResponse(&user{
		ID:        s.fixedUserID,
		Email:     s.fixedEmail,
		FirstName: s.fixedFirstName,
		LastName:  s.fixedLastName,
	})
}

func (s *ServiceTestSuite) TestGetUser() {
	userEntity := &domain.User{}
	userEntity.Write(domain.UserID(s.fixedUserID), s.fixedEmail, "", s.fixedFirstName, s.fixedLastName)
	s.mockUserRepository.On("GetByID", mock.Anything, domain.UserID(s.fixedUserID)).
		Return(userEntity, nil)

	s.serve(http.MethodGet, fmt.Sprint("/v1/users/", s.fixedUserID), nil)

	s.Equal(200, s.responseRecorder.Code)
	s.equalResponse(&user{
		ID:        s.fixedUserID,
		Email:     s.fixedEmail,
		FirstName: s.fixedFirstName,
		LastName:  s.fixedLastName,
	})
}

func (s *ServiceTestSuite) SetupTest() {
	s.container = restful.NewContainer()
	s.mockUserRepository = new(domainmocks.UserRepository)
	s.mockPasswordHashService = new(mocks.PasswordHashService)
	s.service = NewService(s.mockUserRepository, s.mockPasswordHashService)
	s.service.Register(s.container)

	s.responseRecorder = httptest.NewRecorder()

	s.fixedUserID = "SOME_RANDOM_UUID"
	s.fixedEmail = "johndoe@gmail.com"
	s.fixedPassword = "FIXED_PASSWORD"
	s.fixedFirstName = "John"
	s.fixedLastName = "Doe"
	s.fixedHashedPassword = "FIXED_HASHED_PASSWORD"
}

func (s *ServiceTestSuite) serve(method string, path string, requestBody interface{}) {
	jsonMarshaled, err := json.Marshal(&requestBody)
	s.NoError(err)
	request, err := http.NewRequest(method, fmt.Sprint("http://example.com", path), bytes.NewBuffer(jsonMarshaled))
	s.NoError(err)
	request.Header.Add("Content-Type", "application/json")
	s.NoError(err)
	s.container.ServeHTTP(s.responseRecorder, request)
}

func (s *ServiceTestSuite) equalResponse(expect interface{}) {
	actual := reflect.New(reflect.ValueOf(expect).Elem().Type()).Interface()
	err := json.NewDecoder(s.responseRecorder.Body).Decode(actual)
	s.NoError(err)
	s.Equal(expect, actual)
}
