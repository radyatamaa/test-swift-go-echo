package product

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandProduct, userId string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.ProductWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandProduct, userId string) (*models.NewCommandProduct, error)
	GetById(ctx context.Context, id string, token string) (*models.ProductDto, error)
}
