package usecase

import (
	"context"
	"github.com/helper"
	"math"
	"time"

	guuid "github.com/google/uuid"

	"github.com/auth/user"
	"github.com/master/inbound"
	"github.com/models"
)

type inboundUsecase struct {
	userUsecase    user.Usecase
	inboundRepo    inbound.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewinboundUsecase(userUsecase user.Usecase, inboundRepo inbound.Repository, timeout time.Duration) inbound.Usecase {
	return &inboundUsecase{
		userUsecase:    userUsecase,
		inboundRepo:    inboundRepo,
		contextTimeout: timeout,
	}
}
func (m inboundUsecase) Delete(c context.Context, id string, userId string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	err := m.inboundRepo.Delete(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	result := &models.ResponseDelete{
		Id:      id,
		Message: "Success Delete",
	}

	return result, nil
}

func (m inboundUsecase) Update(c context.Context, ar *models.NewCommandInbound, userId string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	getInbound, err := m.inboundRepo.GetByID(ctx, ar.Id)
	if err != nil {
		return err
	}
	var total float64
	total = ar.HargaBeli * float64(ar.Jumlah)
	var modifyBy string = userId
	now := time.Now()

	ib := models.Inbound{
		Id:           getInbound.Id,
		CreatedBy:    getInbound.CreatedBy,
		CreatedDate:  getInbound.CreatedDate,
		ModifiedBy:   &modifyBy,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		InboundTime:  ar.InboundTime,
		ExpiredDate:  ar.ExpiredDate,
		ProductId:    ar.ProductId,
		Jumlah:       ar.Jumlah,
		HargaBeli:     ar.HargaBeli,
		Total:        total,
		NoPO:         getInbound.NoPO,
	}


	err = m.inboundRepo.Update(ctx, &ib)
	if err != nil {
		return err
	}
	return nil
}

func (m inboundUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.InboundWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.inboundRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.InboundDto, len(list))
	for i, item := range list {
		users[i] = &models.InboundDto{
			Id:        item.Id,
			InboundTime 	:        item.InboundTime,
			ProductId 		:        item.ProductId,
			Jumlah 				:        item.Jumlah,
			HargaBeli 		:        item.HargaBeli,
			Total 			:        item.Total,
			NoPO 			:        item.NoPO,
			ProductSKU:item.ProductSKU,
			ProductName:item.ProductName,
		}
	}
	totalRecords, _ := m.inboundRepo.Count(ctx)
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

	response := &models.InboundWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m inboundUsecase) Create(c context.Context, ar *models.NewCommandInbound, userId string) (*models.NewCommandInbound, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	noPO := "ID-" + time.Now().Format("20060102") + "-"
	g, er := helper.GenerateRandomString(5)
	if er != nil {
		return nil, er
	}
	noPO = noPO + g
	var total float64
	total = ar.HargaBeli * float64(ar.Jumlah)
	insert := models.Inbound{
		Id:           guuid.New().String(),
		CreatedBy:    userId,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		InboundTime 	:        ar.InboundTime,
		ProductId 		:        ar.ProductId,
		Jumlah 				:        ar.Jumlah,
		HargaBeli 		:        ar.HargaBeli,
		Total 			:       total,
		NoPO 			:        noPO,
	}

	err := m.inboundRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m inboundUsecase) GetById(c context.Context, id string, userId string) (*models.InboundDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	inbound, err := m.inboundRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.InboundDto{
		Id:        inbound.Id,
		InboundTime 	:        inbound.InboundTime,
		ProductId 		:        inbound.ProductId,
		Jumlah 				:        inbound.Jumlah,
		HargaBeli 		:        inbound.HargaBeli,
		Total 			:        inbound.Total,
		NoPO 			:        inbound.NoPO,
		ProductSKU:inbound.ProductSKU,
		ProductName:inbound.ProductName,
	}

	return result, nil
}
