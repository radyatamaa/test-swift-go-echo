package models

import "time"

type Inbound struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	InboundTime 		time.Time `json:"inbound_time"`
	ExpiredDate 		string `json:"expired_date"`
	ProductId 			string `json:"product_id"`
	Jumlah 				int `json:"jumlah"`
	HargaBeli 			float64 `json:"harga_beli"`
	Total 				float64 `json:"total"`
	NoPO 				string `json:"no_po"`
}

type InboundJoinProduct struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	InboundTime 		time.Time `json:"inbound_time"`
	ExpiredDate 		string `json:"expired_date"`
	ProductId 			string `json:"product_id"`
	Jumlah 				int `json:"jumlah"`
	HargaBeli 			float64 `json:"harga_beli"`
	Total 				float64 `json:"total"`
	NoPO 				string `json:"no_po"`
	ProductSKU 			string `json:"product_sku"`
	ProductName 		string `json:"product_name"`
}


type NewCommandInbound struct {
	Id                   string     `json:"id" validate:"required"`
	InboundTime 		time.Time `json:"inbound_time"`
	ExpiredDate 		string `json:"expired_date"`
	ProductId 			string `json:"product_id"`
	Jumlah 				int `json:"jumlah"`
	HargaBeli 			float64 `json:"harga_beli"`
}

type InboundDto struct {
	Id                   string     `json:"id" validate:"required"`
	ProductSKU 			string `json:"product_sku"`
	ProductName 		string `json:"product_name"`
	InboundTime 		time.Time `json:"inbound_time"`
	ExpiredDate 		string `json:"expired_date"`
	ProductId 			string `json:"product_id"`
	Jumlah 				int `json:"jumlah"`
	HargaBeli 			float64 `json:"harga_beli"`
	Total 				float64 `json:"total"`
	NoPO 				string `json:"no_po"`
}

type InboundWithPagination struct {
	Data []*InboundDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}
