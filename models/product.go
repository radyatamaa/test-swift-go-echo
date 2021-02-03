package models

import "time"

type Product struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	SKU 				string `json:"sku"`
	Name 				string `json:"name"`
	Expirable 			int `json:"expirable"`
}

type NewCommandProduct struct {
	Id                   string     `json:"id" validate:"required"`
	Name 				string `json:"name"`
	Expirable 			int `json:"expirable"`
}

type ProductDto struct {
	Id                   string     `json:"id" validate:"required"`
	SKU 				string `json:"sku"`
	Name 				string `json:"name"`
	Expirable 			int `json:"expirable"`
}

type ProductWithPagination struct {
	Data []*ProductDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}