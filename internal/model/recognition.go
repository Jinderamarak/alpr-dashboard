package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Recognition struct {
	ID        string
	CreatedAt time.Time
	CarID     *string
	Car       *Car `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (recognition *Recognition) BeforeCreate(tx *gorm.DB) (err error) {
	recognition.ID = uuid.NewString()
	return nil
}
