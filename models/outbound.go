package models

import "time"

type Outbound struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	TimeStamp 			time.Time `json:"time_stamp"`
	ProductId 			string `json:"product_id"`
	Qty 				int `json:"qty"`
	Price 				float64 `json:"price"`
	Total 				float64 `json:"total"`
	Usecase 			string `json:"usecase"`
	ReferenceNumber 	string `json:"reference_number"`
}
type OutboundJoinProduct struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	TimeStamp 			time.Time `json:"time_stamp"`
	ProductId 			string `json:"product_id"`
	Qty 				int `json:"qty"`
	Price 				float64 `json:"price"`
	Total 				float64 `json:"total"`
	Usecase 			string `json:"usecase"`
	ReferenceNumber 	string `json:"reference_number"`
	ProductSKU 		string `json:"product_sku"`
	ProductName string `json:"product_name"`
}

type NewCommandOutbound struct {
	Id                   string     `json:"id" validate:"required"`
	TimeStamp 			time.Time `json:"time_stamp"`
	ProductId 			string `json:"product_id"`
	Total 				float64 `json:"total"`
	Usecase 			string `json:"usecase"`
	ReferenceNumber 	string `json:"reference_number"`
	Qty 				float64 `json:"qty"`
	Price 				float64 `json:"price"`
}

type OutboundDto struct {
	Id                   string     `json:"id" validate:"required"`
	TimeStamp 			time.Time `json:"time_stamp"`
	ProductId 			string `json:"product_id"`
	Total 				float64 `json:"total"`
	Usecase 			string `json:"usecase"`
	ReferenceNumber 	string `json:"reference_number"`
	Qty 				float64 `json:"qty"`
	Price 				float64 `json:"price"`
	ProductSKU 		string `json:"product_sku"`
	ProductName string `json:"product_name"`
}

type OutboundWithPagination struct {
	Data []*OutboundDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}
