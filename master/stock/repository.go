package stock

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	DeleteInOutBound(ctx context.Context, inboundId string,outboundId string, deleted_by string) error
	Insert(ctx context.Context, a *models.Stock) error
	Count(ctx context.Context,productId string,bound int) (int, error)
	List(ctx context.Context, limit, offset int,productId string,bound int) ([]*models.StockJoinProductInOutbound, error)
	GetFirst(ctx context.Context, productId string) (*models.Stock, error)
}
