package bcrypt

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHashService struct{}

func NewPasswordHashService() *PasswordHashService {
	return &PasswordHashService{}
}

func (s *PasswordHashService) HashPassword(ctx context.Context, plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrapf(err, "error in generating a hash")
	}
	return string(hashed), nil
}
