package inbound

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.InboundJoinProduct, error)
	Update(ctx context.Context, ar *models.Inbound) error
	Insert(ctx context.Context, a *models.Inbound) error
	Delete(ctx context.Context, id string, deleted_by string) error
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int) ([]*models.InboundJoinProduct, error)
}
