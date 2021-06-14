package domain

import (
	"context"
	"time"
)

type UserID string

type User struct {
	id             UserID
	email          string
	hashedPassword string
	firstName      string
	lastName       string
}

func (u *User) Write(id UserID, email string, hashedPassword string, firstName string, lastName string) {
	u.id = id
	u.email = email
	u.hashedPassword = hashedPassword
	u.firstName = firstName
	u.lastName = lastName
}

func (u *User) Read(f func(id UserID, email string, hashedPassword string, firstName string, lastName string)) {
	f(u.id, u.email, u.hashedPassword, u.firstName, u.lastName)
}

type UserSignedUpEvent struct {
	UserID    UserID    `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *UserSignedUpEvent) Name() string {
	return "UserSignedUpEvent"
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
}
