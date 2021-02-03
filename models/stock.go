package models

import "time"

type Stock struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	ProductId 			string `json:"product_id"`
	InboundId 			*string `json:"inbound_id"`
	OutboundId 			*string `json:"outbound_id"`
	CurrentStock 		int `json:"current_stock"`
}

type StockJoinProductInOutbound struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	ProductId 			string `json:"product_id"`
	InboundId 			*string `json:"inbound_id"`
	OutboundId 			*string `json:"outbound_id"`
	CurrentStock 		int `json:"current_stock"`
	ProductSKU 			string `json:"product_sku"`
	ProductName 		string `json:"product_name"`
	InboundDate			*time.Time `json:"inbound_date"`
	InboundQTY 			*int `json:"inbound_qty"`
	OutboundDate			*time.Time `json:"outbound_date"`
	OutboundQTY 			*int `json:"outbound_qty"`
}

type StockDto struct {
	Id                   string     `json:"id" validate:"required"`
	ProductId 			string `json:"product_id"`
	InboundId 			*string `json:"inbound_id"`
	OutboundId 			*string `json:"outbound_id"`
	CurrentStock 		int `json:"current_stock"`
	ProductSKU 			string `json:"product_sku"`
	ProductName 		string `json:"product_name"`
	InboundDate			*time.Time `json:"inbound_date"`
	InboundQTY 			*int `json:"inbound_qty"`
	OutboundDate			*time.Time `json:"outbound_date"`
	OutboundQTY 			*int `json:"outbound_qty"`
}

type StockWithPagination struct {
	Data []*StockDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}