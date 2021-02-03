package inbound

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context, id string, userId string) (*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandInbound, userId string) error
	List(ctx context.Context, page, limit, offset int, search string) (*models.InboundWithPagination, error)
	Create(ctx context.Context, ar *models.NewCommandInbound, userId string) (*models.NewCommandInbound, error)
	GetById(ctx context.Context, id string, userId string) (*models.InboundDto, error)
}
