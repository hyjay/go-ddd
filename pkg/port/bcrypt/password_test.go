package bcrypt

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PasswordHashServiceTestSuite struct {
	suite.Suite

	passwordHashService *PasswordHashService
}

func TestPasswordHashServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordHashServiceTestSuite))
}

func (s *PasswordHashServiceTestSuite) TestHashPassword() {
	hashed, err := s.passwordHashService.HashPassword(context.Background(), "helloworld")
	s.NoError(err)

	s.True(strings.HasPrefix(hashed, "$2a$10$"))
}

func (s *PasswordHashServiceTestSuite) SetupTest() {
	s.passwordHashService = NewPasswordHashService()
}
