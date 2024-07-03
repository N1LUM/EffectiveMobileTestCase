package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Tasks struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	Status      bool           `gorm:"default:false" json:"status"`
	Hours       int            `gorm:"default:0;not null" json:"hours"`
	Minutes     int            `gorm:"default:0;not null" json:"minutes"`
	Seconds     int            `gorm:"default:0;not null" json:"seconds"`
	StartTime   *time.Time     `gorm:"index"`
	EndTime     *time.Time     `gorm:"index"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt,omitempty"`
}

type TasksCreateInput struct {
	UserID uuid.UUID `json:"user_id"`
	Tasks  Tasks     `json:"task"`
}
