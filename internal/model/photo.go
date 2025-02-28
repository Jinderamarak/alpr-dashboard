package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	ID            string
	RecognitionID *string
	Recognition   *Recognition `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (photo *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	photo.ID = uuid.NewString()
	return nil
}
