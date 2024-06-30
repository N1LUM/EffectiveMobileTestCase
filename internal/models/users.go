package models

import (
	"github.com/google/uuid"
	"time"
)

type Users struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name           string    `gorm:"not null" json:"name"`
	Surname        string    `gorm:"not null" json:"surname"`
	Patronymic     string    `gorm:"not null" json:"patronymic"`
	Address        string    `gorm:"not null" json:"address"`
	PassportSerie  string    `gorm:"not null" json:"passportSerie"`
	PassportNumber string    `gorm:"not null" json:"passportNumber"`
	FullPassport   string    `gorm:"unique"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	DeletedAt      time.Time
}
