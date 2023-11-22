package webapi

import (
	"time"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
)

type WebAPI interface {
	GetUserID(*entity.User) error
	AddAbsenceStatus(*entity.User, [2]time.Time) error
}
