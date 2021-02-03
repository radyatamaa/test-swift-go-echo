package order

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Insert(ctx context.Context, a *models.Order) error
	Count(ctx context.Context,referenceId string) (int, error)
	List(ctx context.Context, limit, offset int,referenceId string) ([]*models.Order, error)
}
