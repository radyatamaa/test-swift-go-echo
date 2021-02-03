package usecase

import (
	"context"
	"github.com/helper"
	"github.com/master/outbound"
	"github.com/master/stock"
	"math"
	"time"

	guuid "github.com/google/uuid"

	"github.com/auth/user"
	"github.com/models"
	"github.com/order/order"
)

type orderUsecase struct {
	outboundRepo outbound.Repository
	userUsecase    user.Usecase
	orderRepo      order.Repository
	stockRepo 		stock.Repository
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NeworderUsecase(userUsecase user.Usecase, orderRepo order.Repository, timeout time.Duration,
	outboundRepo outbound.Repository,	stockRepo 		stock.Repository) order.Usecase {
	return &orderUsecase{
		stockRepo:stockRepo,
		outboundRepo:outboundRepo,
		userUsecase:    userUsecase,
		orderRepo:      orderRepo,
		contextTimeout: timeout,
	}
}

func (m orderUsecase) UpdateStatus(c context.Context, order *models.NewCommandOrderStatus) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	orders, err := m.orderRepo.List(ctx, 1,0,order.ReferenceNumber)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		return models.ErrNotFound
	}

	getOrder := orders[0]

	now := time.Now()
	getOrder.CreatedDate = now
	getOrder.Id = guuid.New().String()
	if order.Status == "pending"{
		getOrder.Status = 0
	}else if order.Status == "picking"{
		getOrder.Status = 1
	}else if order.Status == "ready_to_pack"{
		getOrder.Status = 2
	}else if order.Status == "packing"{
		getOrder.Status = 3
	}else if order.Status == "routing"{
		getOrder.Status = 4
	}else if order.Status == "shipped"{
		getOrder.Status = 5
	}else if order.Status == "delivered"{
		getOrder.Status = 6
	}else if order.Status == "ftd"{
		getOrder.Status = 7
	}else if order.Status == "cancelled"{
		getOrder.Status = 8
	}

	err = m.orderRepo.Insert(ctx, getOrder)


	if err != nil {
		return err
	}
	return nil
}

func (m orderUsecase) List(ctx context.Context, page, limit, offset int, referenceNumber string) (*models.OrderWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.orderRepo.List(ctx, limit, offset,referenceNumber)
	if err != nil {
		return nil, err
	}

	users := make([]*models.OrderDto, len(list))
	for i, item := range list {
		var status string
		if item.Status == 0{
			status = "pending"
		}else if item.Status == 1{
			status = "picking"
		}else if item.Status == 2{
			status = "ready_to_pack"
		}else if item.Status == 3{
			status = "packing"
		}else if item.Status == 4{
			status = "routing"
		}else if item.Status == 5{
			status = "shipped"
		}else if item.Status == 6{
			status = "delivered"
		}else if item.Status == 7{
			status = "ftd"
		}else if item.Status == 8{
			status = "cancelled"
		}
		users[i] = &models.OrderDto{
			Id:        item.Id,
			Date:item.CreatedDate.Format("2006-01-02"),
			ReferenceNumber:  item.ReferenceNumber,
			CustomerName:    item.ReferenceNumber,
			SourceAddress:    item.SourceAddress,
			DestAddress:      item.DestAddress,
			Status:            status,
			TotalPrice:        item.TotalPrice,
			CustomerReceived:  item.CustomerReceived,
			Remarks:           item.Remarks,
		}
	}
	totalRecords, _ := m.orderRepo.Count(ctx,referenceNumber)
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

	response := &models.OrderWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m orderUsecase) Create(c context.Context, ar *models.NewCommandOrder, userId string) (*models.NewCommandOrder, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	referenceNumber := "ID-" + time.Now().Format("20060102") + "-"
	g, er := helper.GenerateRandomString(5)
	if er != nil {
		return nil, er
	}

	referenceNumber = referenceNumber + g
	var totalPrice float64 = 0
	for _,element := range ar.Product{
		total := element.Price * float64(element.Qty)
		ouboundM := models.Outbound{
			Id:              guuid.New().String(),
			CreatedBy:       userId,
			CreatedDate:     time.Now(),
			ModifiedBy:      nil,
			ModifiedDate:    nil,
			DeletedBy:       nil,
			DeletedDate:     nil,
			IsDeleted:       0,
			IsActive:        1,
			TimeStamp:       time.Now(),
			ProductId:       element.ProductId,
			Qty:             element.Qty,
			Price:           element.Price,
			Total:           total,
			Usecase:         element.Usecase,
			ReferenceNumber: referenceNumber,
		}

		err := m.outboundRepo.Insert(ctx, &ouboundM)
		if err != nil {
			return nil, err
		}



		getStock ,err := m.stockRepo.GetFirst(ctx,element.ProductId)
		if err != nil {
			return nil, err
		}

		var currentStock int
		if getStock != nil {
			currentStock = getStock.CurrentStock - element.Qty
		}

		stock := models.Stock{
			Id:           guuid.New().String(),
			CreatedBy:    userId,
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ProductId:    element.ProductId,
			InboundId:    nil,
			OutboundId:   &ouboundM.Id,
			CurrentStock: currentStock,
		}

		err = m.stockRepo.Insert(ctx,&stock)
		if err != nil {
			return nil, err
		}

		totalPrice = totalPrice + total

	}
	insert := models.Order{
		Id:           guuid.New().String(),
		CreatedBy:    userId,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		ReferenceNumber:  referenceNumber,
		CustomerName:     ar.CustomerName,
		SourceAddress:    ar.SourceAddress,
		DestAddress:      ar.DestAddress,
		Status:           0,
		TotalPrice:       totalPrice,
		CustomerReceived: ar.CustomerReceived,
		Remarks:          ar.Remarks,
	}

	err := m.orderRepo.Insert(ctx, &insert)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

