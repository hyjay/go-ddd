package fakes

import (
	"context"
	"github.com/hyjay/go-ddd/pkg/domain"
	"github.com/pkg/errors"
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

func (r *UserRepository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	user, ok := r.users[string(id)]
	if !ok {
		return nil, errors.New("no such user")
	}
	userCopy := *user
	return &userCopy, nil
}
