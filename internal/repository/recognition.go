package repository

import (
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"github.com/jinderamarak/alpr-dasboard/internal/util"
	"gorm.io/gorm"
)

type RecognitionRepository interface {
	GetByIdWithCar(id uuid.UUID) (model.Recognition, error)
	GetPageWithCarId(carId *uuid.UUID, offset, limit int) ([]model.Recognition, error)
	CountWithCarId(carId *uuid.UUID) (int64, error)
	GetPageWithCar(offset, limit int) ([]model.Recognition, error)
	Count() (int64, error)
	Create(recognition model.Recognition) (model.Recognition, error)
}

type recognitionRepository struct {
	db *gorm.DB
}

func NewRecognitionRepository(db *gorm.DB) RecognitionRepository {
	return &recognitionRepository{db}
}

func (repo *recognitionRepository) GetByIdWithCar(id uuid.UUID) (model.Recognition, error) {
	recognition := model.Recognition{ID: id.String()}
	result := repo.db.Joins("Car").First(&recognition)
	return recognition, result.Error
}

func (repo *recognitionRepository) GetPageWithCarId(carId *uuid.UUID, offset, limit int) ([]model.Recognition, error) {
	id := util.MapPtr(carId, func(uuid uuid.UUID) string {
		return uuid.String()
	})

	var recognitions []model.Recognition
	result := repo.db.Find(&recognitions).Order("created_at desc").Where("car_id = ?", id).Offset(offset).Limit(limit)
	return recognitions, result.Error
}

func (repo *recognitionRepository) CountWithCarId(carId *uuid.UUID) (int64, error) {
	id := util.MapPtr(carId, func(uuid uuid.UUID) string {
		return uuid.String()
	})

	var count int64
	result := repo.db.Find(&model.Recognition{}).Where("car_id = ?", id).Count(&count)
	return count, result.Error
}

func (repo *recognitionRepository) GetPageWithCar(offset, limit int) ([]model.Recognition, error) {
	var recognitions []model.Recognition
	result := repo.db.Joins("Car").Order("created_at desc").Offset(offset).Limit(limit).Find(&recognitions)
	return recognitions, result.Error
}

func (repo *recognitionRepository) Count() (int64, error) {
	var count int64
	result := repo.db.Find(&model.Recognition{}).Count(&count)
	return count, result.Error
}

func (repo *recognitionRepository) Create(recognition model.Recognition) (model.Recognition, error) {
	result := repo.db.Create(&recognition)
	return recognition, result.Error
}
