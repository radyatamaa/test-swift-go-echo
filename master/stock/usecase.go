package stock

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context, page, limit, offset int, productId string,bound int) (*models.StockWithPagination, error)
}
