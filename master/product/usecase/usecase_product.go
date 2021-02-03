package usecase

import (
	"context"
	guuid "github.com/google/uuid"
	"github.com/helper"
	"math"
	"time"

	"github.com/auth/user"
	"github.com/master/product"
	"github.com/models"
)

type productUsecase struct {
	userUsecase    user.Usecase
	productRepo    product.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewproductUsecase(userUsecase user.Usecase, productRepo product.Repository, timeout time.Duration) product.Usecase {
	return &productUsecase{
		userUsecase:    userUsecase,
		productRepo:    productRepo,
		contextTimeout: timeout,
	}
}
func (m productUsecase) Delete(c context.Context, id string, userId string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	err := m.productRepo.Delete(ctx, id, userId)
	if err != nil {
		return nil,err
	}
	result := &models.ResponseDelete{
		Id:    id,
		Message: "Success Delete",
	}

	return result, nil
}

func (m productUsecase) Update(c context.Context, ar *models.NewCommandProduct, userId string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	getProduct, err := m.productRepo.GetByID(ctx, ar.Id)
	if err != nil {
		return err
	}
	var modifyBy string = userId
	now := time.Now()

	getProduct.Name  =       ar.Name
	getProduct.Expirable =   ar.Expirable
	getProduct.ModifiedBy = &modifyBy
	getProduct.ModifiedDate = &now
	err = m.productRepo.Update(ctx, getProduct)
	if err != nil {
		return err
	}
	return nil
}

func (m productUsecase) List(ctx context.Context, page, limit, offset int, search string) (*models.ProductWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.productRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.ProductDto, len(list))
	for i, item := range list {
		users[i] = &models.ProductDto{
			Id:        item.Id,
			SKU:       item.SKU,
			Name:      item.Name,
			Expirable: item.Expirable,
		}
	}
	totalRecords, _ := m.productRepo.Count(ctx)
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

	response := &models.ProductWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m productUsecase) Create(c context.Context, ar *models.NewCommandProduct, userId string) (*models.NewCommandProduct, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	g, er := helper.GenerateRandomStringWithChar(20)
	if er != nil {
		return nil, er
	}
	insert := models.Product{
		Id:           guuid.New().String(),
		CreatedBy:    userId,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		SKU:          g,
		Name:         ar.Name,
		Expirable:    ar.Expirable,
	}

	err := m.productRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func (m productUsecase) GetById(c context.Context, id string, userId string) (*models.ProductDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	product, err := m.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	result := &models.ProductDto{
		Id:        product.Id,
		SKU:       product.SKU,
		Name:      product.Name,
		Expirable: product.Expirable,
	}

	return result, nil
}
