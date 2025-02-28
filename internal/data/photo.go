package data

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"github.com/jinderamarak/alpr-dasboard/internal/util"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"net/url"
	"time"
)

const (
	bucket    = "alpr"
	keyPrefix = "photo/"
)

var (
	createPresignedExpiration = time.Duration(10) * time.Minute
	getPresignedExpiration    = time.Duration(10) * time.Minute
)

type PhotoRepository interface {
	CreateUploadByRecognitionId(recognitionId uuid.UUID) (*url.URL, map[string]string, error)
	GetByRecognitionId(recognitionId *uuid.UUID) ([]url.URL, error)
}

type photoRepository struct {
	db *gorm.DB
	s3 *minio.Client
}

func NewPhotoRepository(db *gorm.DB, s3 *minio.Client) PhotoRepository {
	if db == nil {
		panic("expected db, but got nil")
	}
	if s3 == nil {
		panic("expected s3, but got nil")
	}
	return &photoRepository{db, s3}
}

func (repo *photoRepository) CreateUploadByRecognitionId(recognitionId uuid.UUID) (*url.URL, map[string]string, error) {
	id := recognitionId.String()
	photo := model.Photo{RecognitionID: &id}

	result := repo.db.Create(&photo)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	policy := minio.NewPostPolicy()
	errors.Unwrap(policy.SetBucket(bucket))
	errors.Unwrap(policy.SetKey(keyPrefix + photo.ID))
	errors.Unwrap(policy.SetExpires(time.Now().UTC().Add(createPresignedExpiration)))
	errors.Unwrap(policy.SetContentType("image/jpeg"))

	return repo.s3.PresignedPostPolicy(context.Background(), policy)
}

func (repo *photoRepository) GetByRecognitionId(recognitionId *uuid.UUID) ([]url.URL, error) {
	id := util.MapPtr(recognitionId, func(uuid uuid.UUID) string {
		return uuid.String()
	})

	var photos []model.Photo
	result := repo.db.Where("recognition_id = ?", id).Find(&photos)
	if result.Error != nil {
		return nil, result.Error
	}

	urls := make([]url.URL, len(photos))
	reqParams := make(url.Values)
	for i, p := range photos {
		presigned, err := repo.s3.PresignedGetObject(context.Background(), bucket, keyPrefix+p.ID, getPresignedExpiration, reqParams)
		if err != nil {
			return nil, err
		}
		urls[i] = *presigned
	}

	return urls, nil
}
