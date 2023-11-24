package usecase

import (
	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/AnatoliyBr/data-modifier/internal/webapi"
)

type AppUseCase struct {
	UserWebAPI webapi.WebAPI
}

func NewAppUseCase(a webapi.WebAPI) *AppUseCase {
	return &AppUseCase{
		UserWebAPI: a,
	}
}

func (uc *AppUseCase) GetUserID(u *entity.User) error {
	return uc.UserWebAPI.GetUserID(u)
}

func (uc *AppUseCase) AddAbsenceStatus(u *entity.User, p [2]entity.CustomTime) error {
	return uc.UserWebAPI.AddAbsenceStatus(u, p)
}
