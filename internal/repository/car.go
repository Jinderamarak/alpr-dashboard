package repository

import (
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"gorm.io/gorm"
)

type CarRepository interface {
	GetOrCreateByPlate(plate string) (model.Car, error)
	GetById(id uuid.UUID) (model.Car, error)
	GetPage(offset, limit int) ([]model.Car, error)
	Count() (int64, error)
}

type carRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) CarRepository {
	return &carRepository{db}
}

func (repo *carRepository) GetOrCreateByPlate(plate string) (model.Car, error) {
	var car model.Car
	result := repo.db.Where(model.Car{Plate: plate}).FirstOrCreate(&car)
	return car, result.Error
}

func (repo *carRepository) GetById(id uuid.UUID) (model.Car, error) {
	car := model.Car{ID: id.String()}
	result := repo.db.First(&car)
	return car, result.Error
}

func (repo *carRepository) GetPage(offset, limit int) ([]model.Car, error) {
	var cars []model.Car
	result := repo.db.Order("created_at desc").Offset(offset).Limit(limit).Find(&cars)
	return cars, result.Error
}

func (repo *carRepository) Count() (int64, error) {
	var count int64
	result := repo.db.Find(&model.Car{}).Count(&count)
	return count, result.Error
}
