package user

import (
	"context"
	"github.com/models"
)

type Repository interface {
	ValidateUser(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user models.User) (*string, error)
}
