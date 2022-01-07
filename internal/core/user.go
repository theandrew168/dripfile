package core

import (
	"context"
)

type User struct {
	Name     string
	Email    string
	Password string
	Verified bool
	Role     string
	Account  Account

	// readonly (from database, after creation)
	ID int
}

func NewUser(name, email, password string, account Account) User {
	user := User{
		Name:     name,
		Email:    email,
		Password: password,
		Verified: false,
		Role:     "viewer",
		Account:  account,
	}
	return user
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *User) error
	ReadUser(ctx context.Context, id int) (User, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, user User) error
}
