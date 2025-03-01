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

type RecognitionImageRepository interface {
	CreateUploadByRecognitionId(recognitionId uuid.UUID) (*url.URL, map[string]string, error)
	ListByRecognitionId(recognitionId *uuid.UUID) ([]*url.URL, error)
}

type imageRepository struct {
	db      *gorm.DB
	objects *minio.Client
}

func NewRecognitionImageRepository(db *gorm.DB, objects *minio.Client) RecognitionImageRepository {
	if db == nil {
		panic("expected db, but got nil")
	}
	if objects == nil {
		panic("expected objects, but got nil")
	}
	return &imageRepository{db, objects}
}

func (repo *imageRepository) CreateUploadByRecognitionId(recognitionId uuid.UUID) (*url.URL, map[string]string, error) {
	id := recognitionId.String()
	photo := model.RecognitionImage{RecognitionID: &id}

	result := repo.db.Create(&photo)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	policy := minio.NewPostPolicy()
	errors.Unwrap(policy.SetBucket(bucket))
	errors.Unwrap(policy.SetKey(keyPrefix + photo.ID))
	errors.Unwrap(policy.SetExpires(time.Now().UTC().Add(createPresignedExpiration)))
	errors.Unwrap(policy.SetContentType("image/jpeg"))

	return repo.objects.PresignedPostPolicy(context.Background(), policy)
}

func (repo *imageRepository) ListByRecognitionId(recognitionId *uuid.UUID) ([]*url.URL, error) {
	id := util.MapPtr(recognitionId, func(uuid uuid.UUID) string {
		return uuid.String()
	})

	var photos []model.RecognitionImage
	result := repo.db.Where("recognition_id = ?", id).Find(&photos)
	if result.Error != nil {
		return nil, result.Error
	}

	urls := make([]*url.URL, len(photos))
	reqParams := make(url.Values)
	for i, p := range photos {
		presigned, err := repo.objects.PresignedGetObject(context.Background(), bucket, keyPrefix+p.ID, getPresignedExpiration, reqParams)
		if err != nil {
			return nil, err
		}
		urls[i] = presigned
	}

	return urls, nil
}
