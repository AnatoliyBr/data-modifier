package usecase

import (
	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/AnatoliyBr/data-modifier/internal/webapi"
)

// AppUseCase implements UseCase interface. It implements
// data-modifier business logic through third-party WebAPI calls.
type AppUseCase struct {
	UserWebAPI webapi.WebAPI
}

// NewAppUseCase returns AppUseCase.
func NewAppUseCase(a webapi.WebAPI) *AppUseCase {
	return &AppUseCase{
		UserWebAPI: a,
	}
}

// GetUserID finds the user id on a WebAPI.
func (uc *AppUseCase) GetUserID(u *entity.User) error {
	return uc.UserWebAPI.GetUserID(u)
}

// AddAbsenceStatus finds user absence status info
// for a time period p on a WebAPI and adds it to the user data.
func (uc *AppUseCase) AddAbsenceStatus(u *entity.User, p [2]entity.CustomTime) error {
	return uc.UserWebAPI.AddAbsenceStatus(u, p)
}
