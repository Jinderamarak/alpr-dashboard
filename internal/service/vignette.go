package service

import (
	"github.com/jinderamarak/alpr-dasboard/internal/data"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"time"
)

type VignetteService interface {
	ValidatePlate(plate string) (model.VignetteResult, error)
}

const expirationSoon = time.Second * time.Duration(10)

type vignetteService struct {
	token    model.VignetteAuth
	provider data.VignetteProvider
}

func NewVignetteService(provider data.VignetteProvider) VignetteService {
	return &vignetteService{
		token:    model.VignetteAuth{},
		provider: provider,
	}
}

func (service *vignetteService) refreshToken() error {
	auth, err := service.provider.GetAuthToken()
	if err != nil {
		return err
	}

	service.token = auth
	return nil
}

func (service *vignetteService) ValidatePlate(plate string) (model.VignetteResult, error) {
	if service.token.ExpiresSoon(expirationSoon) {
		err := service.refreshToken()
		if err != nil {
			return model.VignetteResult{}, err
		}
	}

	return service.provider.GetVignetteStatus(plate, &service.token)
}
