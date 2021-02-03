package models

import (
	"time"
)

type User struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	UserEmail            string     `json:"user_email" validate:"required"`
	Password string `json:"password"`
	Phone string `json:"phone"`
}

type UserDto struct {
	Id                   string     `json:"id" validate:"required"`
	UserEmail            string     `json:"user_email" validate:"required"`
	Password string `json:"password"`
	Phone string `json:"phone"`
}

type NewCommandUser struct {
	Id                   string     `json:"id" validate:"required"`
	UserEmail            string     `json:"user_email" validate:"required"`
	Password string `json:"password"`
	Phone string `json:"phone"`
}

