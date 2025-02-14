package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Car struct {
	ID           string
	Plate        string `gorm:"unique"`
	Description  string
	IsAuthorized bool
	CreatedAt    time.Time
	DeletedAt    *time.Time
}

func (car *Car) BeforeCreate(tx *gorm.DB) (err error) {
	car.ID = uuid.NewString()
	return nil
}
