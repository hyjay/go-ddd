package fakes

import (
	"context"
	"github.com/hyjay/go-ddd/pkg/domain"
	"sync"
)

type UserRepository struct {
	users map[string]*domain.User
	mutex *sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
		mutex: &sync.RWMutex{},
	}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	userCopy := *user
	r.users[string(user.ID())] = &userCopy
	return nil
}
