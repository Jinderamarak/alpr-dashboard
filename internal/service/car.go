package service

import (
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"github.com/jinderamarak/alpr-dasboard/internal/repository"
	"github.com/jinderamarak/alpr-dasboard/internal/util"
	"strings"
)

const CarPageSize = 10

type CarService interface {
	GetPage(page int) ([]model.Car, error)
	CountPages() (int, error)
	GetById(carId uuid.UUID) (model.Car, error)
	GetOrCreateByPlate(plate string) (model.Car, error)
}

type carService struct {
	cars repository.CarRepository
}

func NewCarService(cars repository.CarRepository) CarService {
	return &carService{cars}
}

func (service *carService) GetPage(page int) ([]model.Car, error) {
	offset := (page - 1) * CarPageSize
	limit := CarPageSize
	return service.cars.GetPage(offset, limit)
}

func (service *carService) CountPages() (int, error) {
	count, err := service.cars.Count()
	if err != nil {
		return 0, err
	}
	return util.NumberOfPages(count, CarPageSize), nil
}

func (service *carService) GetById(carId uuid.UUID) (model.Car, error) {
	return service.cars.GetById(carId)
}

func (service *carService) GetOrCreateByPlate(plate string) (model.Car, error) {
	plate = strings.ToUpper(plate)
	return service.cars.GetOrCreateByPlate(plate)
}
