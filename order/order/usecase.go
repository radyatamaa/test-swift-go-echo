package order

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	UpdateStatus(ctx context.Context, order *models.NewCommandOrderStatus) error
	List(ctx context.Context, page, limit, offset int, referenceNumber string) (*models.OrderWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandOrder, userId string) (*models.NewCommandOrder, error)
}
