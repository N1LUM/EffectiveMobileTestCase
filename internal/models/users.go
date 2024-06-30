package models

import (
	"github.com/google/uuid"
	"time"
)

type Users struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name           string    `gorm:"not null"`
	Surname        string    `gorm:"not null"`
	Patronymic     string    `gorm:"not null"`
	Address        string    `gorm:"not null"`
	PassportSerie  string    `gorm:"not null"`
	PassportNumber string    `gorm:"not null"`
	FullPassport   string    `gorm:"unique"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	DeletedAt      time.Time
}
