package model

import "time"

type Category struct {
	Id        int64 `gorm:"primary_key"`
	Title     string
	IsPrivate bool `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CategoryStore interface {
	CategoryList(onlyPrivate bool) ([]*Category, error)
	CategoryCreate(category *Category) error
	CategoryLoad(cId int64) (*Category, error)
	CategoryDelete(cId int64) error
	CategoryUpdate(category *Category) error
}

type CategoryService interface {
	CategoryStore
}
