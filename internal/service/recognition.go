package service

import (
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/data"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"github.com/jinderamarak/alpr-dasboard/internal/util"
)

const RecognitionPageSize = 10

type RecognitionService interface {
	CreateByPlate(plate string) (model.Recognition, error)
	GetPage(page int) ([]model.Recognition, error)
	CountPages() (int, error)
	GetPageWithCarId(carId *uuid.UUID, page int) ([]model.Recognition, error)
	CountPagesWithCarId(carId *uuid.UUID) (int, error)
	GetByIdWithCar(recognitionId uuid.UUID) (model.Recognition, error)
	Notifications() *util.Broker[model.Recognition]
}

type recognitionService struct {
	recognitions  data.RecognitionRepository
	carService    CarService
	notifications *util.Broker[model.Recognition]
}

func NewRecognitionService(recognitions data.RecognitionRepository, carService CarService) RecognitionService {
	notifications := util.NewBroker[model.Recognition]()
	go notifications.Start()
	return &recognitionService{recognitions, carService, notifications}
}

func (service *recognitionService) CreateByPlate(plate string) (model.Recognition, error) {
	car, err := service.carService.GetOrCreateByPlate(plate)
	if err != nil {
		return model.Recognition{}, err
	}

	recognition := model.Recognition{
		CarID: &car.ID,
	}

	created, err := service.recognitions.Create(recognition)
	if err != nil {
		return model.Recognition{}, err
	}

	service.notifications.Publish(created)
	return created, nil
}

func (service *recognitionService) GetPage(page int) ([]model.Recognition, error) {
	offset := (page - 1) * RecognitionPageSize
	limit := RecognitionPageSize
	return service.recognitions.GetPageWithCar(offset, limit)
}

func (service *recognitionService) CountPages() (int, error) {
	count, err := service.recognitions.Count()
	if err != nil {
		return 0, err
	}
	return util.NumberOfPages(count, RecognitionPageSize), nil
}

func (service *recognitionService) GetPageWithCarId(carId *uuid.UUID, page int) ([]model.Recognition, error) {
	offset := (page - 1) * RecognitionPageSize
	limit := RecognitionPageSize
	return service.recognitions.GetPageWithCarId(carId, offset, limit)
}

func (service *recognitionService) CountPagesWithCarId(carId *uuid.UUID) (int, error) {
	count, err := service.recognitions.CountWithCarId(carId)
	if err != nil {
		return 0, err
	}
	return util.NumberOfPages(count, RecognitionPageSize), nil
}

func (service *recognitionService) GetByIdWithCar(recognitionId uuid.UUID) (model.Recognition, error) {
	return service.recognitions.GetByIdWithCar(recognitionId)
}

func (service *recognitionService) Notifications() *util.Broker[model.Recognition] {
	return service.notifications
}
