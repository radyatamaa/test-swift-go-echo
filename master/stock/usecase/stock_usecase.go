package usecase

import (
	"context"
	"math"
	"time"

	"github.com/auth/user"
	"github.com/master/stock"
	"github.com/models"
)

type stockUsecase struct {
	userUsecase    user.Usecase
	stockRepo      stock.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewstockUsecase(userUsecase user.Usecase, stockRepo stock.Repository, timeout time.Duration) stock.Usecase {
	return &stockUsecase{
		userUsecase:    userUsecase,
		stockRepo:      stockRepo,
		contextTimeout: timeout,
	}
}
func (m stockUsecase) List(ctx context.Context, page, limit, offset int,productId string,bound int) (*models.StockWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.stockRepo.List(ctx, limit, offset,productId,bound)
	if err != nil {
		return nil, err
	}

	users := make([]*models.StockDto, len(list))
	for i, item := range list {
		users[i] = &models.StockDto{
			Id:        item.Id,
			ProductId 	:        item.Id,
			InboundId 		:        item.InboundId,
			OutboundId 		:        item.OutboundId,
			CurrentStock 	:        item.CurrentStock,
			ProductSKU 		:        item.ProductSKU,
			ProductName 	:        item.ProductName,
			InboundDate			:        item.InboundDate,
			InboundQTY 			:        item.InboundQTY,
			OutboundDate		:        item.OutboundDate,
			OutboundQTY 			:        item.OutboundQTY,
		}
	}
	totalRecords, _ := m.stockRepo.Count(ctx,productId,bound)
	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}
	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(list),
	}

	response := &models.StockWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}
