package outbound

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.Outbound, error)
	Update(ctx context.Context, ar *models.Outbound) error
	Insert(ctx context.Context, a *models.Outbound) error
	Delete(ctx context.Context, id string, deleted_by string) error
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int) ([]*models.OutboundJoinProduct, error)
}
