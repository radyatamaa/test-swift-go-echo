package outbound

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id string, userId string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandOutbound, userId string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.OutboundWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandOutbound, userId string) (*models.NewCommandOutbound, error)
	GetById(ctx context.Context, id string, token string) (*models.OutboundDto, error)
}
