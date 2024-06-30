package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Tasks struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Status      bool      `gorm:"default:false"`
	Hours       int       `gorm:"default:0;not null"`
	Minutes     int       `gorm:"default:0;not null"`
	Seconds     int       `gorm:"default:0;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
}
