package model

import "time"

type Category struct {
	Id        int64 `gorm:"primary_key"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CategoryStore interface {
}

type CategoryService interface {
	CategoryStore
}
