package usecase

import (
	"context"
	"math"
	"time"

	guuid "github.com/google/uuid"

	"github.com/auth/user"
	"github.com/master/outbound"
	"github.com/models"
)

type outboundUsecase struct {
	userUsecase    user.Usecase
	outboundRepo   outbound.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewoutboundUsecase(userUsecase user.Usecase, outboundRepo outbound.Repository, timeout time.Duration) outbound.Usecase {
	return &outboundUsecase{
		userUsecase:    userUsecase,
		outboundRepo:   outboundRepo,
		contextTimeout: timeout,
	}
}
func (m outboundUsecase) Delete(c context.Context, id string, userId string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	err := m.outboundRepo.Delete(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	result := &models.ResponseDelete{
		Id:      id,
		Message: "Success Delete",
	}

	return result, nil
}

func (m outboundUsecase) Update(c context.Context, ar *models.NewCommandOutbound, userId string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	getOutbound, err := m.outboundRepo.GetByID(ctx, ar.Id)
	if err != nil {
		return err
	}
	var modifyBy string = userId
	now := time.Now()
	getOutbound.TimeStamp 	=           ar.TimeStamp
	getOutbound.ProductId 		=           ar.ProductId
	getOutbound.	Total 			=           ar.Total
	getOutbound.	Usecase 			=           ar.Usecase
	getOutbound.	ReferenceNumber 	=           ar.ReferenceNumber
	getOutbound.ModifiedBy = &modifyBy
	getOutbound.ModifiedDate = &now
	err = m.outboundRepo.Update(ctx, getOutbound)
	if err != nil {
		return err
	}
	return nil
}

func (m outboundUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.OutboundWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.outboundRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.OutboundDto, len(list))
	for i, item := range list {
		users[i] = &models.OutboundDto{
			Id:           item.Id,
			TimeStamp 	:           item.TimeStamp,
			ProductId 		:           item.ProductId,
			Total 			:           item.Total,
			Usecase 			:           item.Usecase,
			ReferenceNumber 	:           item.ReferenceNumber,
			ProductSKU:item.ProductSKU,
			ProductName:item.ProductName,
		}
	}
	totalRecords, _ := m.outboundRepo.Count(ctx)
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

	response := &models.OutboundWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m outboundUsecase) Create(c context.Context, ar *models.NewCommandOutbound, userId string) (*models.NewCommandOutbound, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	insert := models.Outbound{
		Id:           guuid.New().String(),
		CreatedBy:    userId,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		TimeStamp 	:           ar.TimeStamp,
		ProductId 		:           ar.ProductId,
		Total 			:           ar.Total,
		Usecase 			:           ar.Usecase,
		ReferenceNumber 	:           ar.ReferenceNumber,
	}

	err := m.outboundRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m outboundUsecase) GetById(c context.Context, id string, userId string) (*models.OutboundDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	outbound, err := m.outboundRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.OutboundDto{
		Id:           outbound.Id,
		TimeStamp 	:           outbound.TimeStamp,
		ProductId 		:           outbound.ProductId,
		Total 			:           outbound.Total,
		Usecase 			:           outbound.Usecase,
		ReferenceNumber 	:           outbound.ReferenceNumber,
	}

	return result, nil
}
