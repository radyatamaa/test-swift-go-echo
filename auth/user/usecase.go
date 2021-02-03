package user

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	ValidateUser(ctx context.Context, email,password string) (*models.UserDto, error)
	Create(ctx context.Context, user models.NewCommandUser) (*models.ResponseDelete, error)
}