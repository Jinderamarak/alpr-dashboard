package service

import (
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/data"
	"net/url"
)

type PhotoService interface {
	CreateUploadByRecognitionId(recognitionId uuid.UUID) (*url.URL, map[string]string, error)
	ListByRecognitionId(recognitionId *uuid.UUID) ([]*url.URL, error)
}

type photoService struct {
	photos data.PhotoRepository
}

func NewPhotoService(photos data.PhotoRepository) PhotoService {
	if photos == nil {
		panic("expected PhotoRepository, got nil")
	}
	return &photoService{photos}
}

func (service *photoService) CreateUploadByRecognitionId(recognitionId uuid.UUID) (*url.URL, map[string]string, error) {
	return service.photos.CreateUploadByRecognitionId(recognitionId)
}

func (service *photoService) ListByRecognitionId(recognitionId *uuid.UUID) ([]*url.URL, error) {
	return service.photos.ListByRecognitionId(recognitionId)
}
