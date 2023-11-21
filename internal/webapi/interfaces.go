package webapi

import "github.com/AnatoliyBr/data-modifier/internal/entity"

type WebAPI interface {
	Authenticate() error
	AddAbsenceStatus(*entity.User) (*entity.User, error)
}
