package models

import "time"

type Article struct {
	ID        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Title     string `json:"title" validate:"required"`
	Content   string `json:"content" validate:"required"`
	Author    Author `json:"author"`
}
