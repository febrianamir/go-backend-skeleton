package model

import "time"

type AppSetting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uint      `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uint      `json:"updated_by"`
	DeletedAt time.Time `json:"deleted_at"`
	DeletedBy uint      `json:"deleted_by"`
}
