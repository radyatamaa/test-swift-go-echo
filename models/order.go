package models

import "time"

type Order struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	ReferenceNumber string `json:"reference_number"`
	CustomerName 	string `json:"customer_name"`
	SourceAddress 	string `json:"source_address"`
	DestAddress 	string `json:"dest_address"`
	Status 			int `json:"status"`
	TotalPrice 		float64 `json:"total_price"`
	CustomerReceived string `json:"customer_received"`
	Remarks 		string `json:"remarks"`
}

type NewCommandOrder struct {
	Id                   string     `json:"id" validate:"required"`
	CustomerName 	string `json:"customer_name"`
	SourceAddress 	string `json:"source_address"`
	DestAddress 	string `json:"dest_address"`
	CustomerReceived string `json:"customer_received"`
	Remarks 		string `json:"remarks"`
	Product        []ProductOutboundObj `json:"product"`
}

type ProductOutboundObj struct {
	ProductId 			string `json:"product_id"`
	Qty 				int `json:"qty"`
	Price 				float64 `json:"price"`
	Usecase 			string `json:"usecase"`
}
type NewCommandOrderStatus struct {
	ReferenceNumber string `json:"reference_number"`
	Status 			string `json:"status"`
	CustomerReceived string `json:"customer_received"`
	Remarks 		string `json:"remarks"`

}
type OrderDto struct {
	Id                   string     `json:"id" validate:"required"`
	Date 			string `json:"date"`
	ReferenceNumber string `json:"reference_number"`
	CustomerName 	string `json:"customer_name"`
	SourceAddress 	string `json:"source_address"`
	DestAddress 	string `json:"dest_address"`
	Status 			string `json:"status"`
	TotalPrice 		float64 `json:"total_price"`
	CustomerReceived string `json:"customer_received"`
	Remarks 		string `json:"remarks"`
}

type OrderWithPagination struct {
	Data []*OrderDto  `json:"data"`
	Meta *MetaPagination `json:"meta"`
}

