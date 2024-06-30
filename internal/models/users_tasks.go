package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UsersTasks struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID `gorm:"index"`
	TaskID    uuid.UUID `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}
