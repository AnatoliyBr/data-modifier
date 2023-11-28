// Package webapi implements working with a specific WebAPI.
package webapi

import (
	"github.com/AnatoliyBr/data-modifier/internal/entity"
)

// WebAPI is a minimal interface for modifying
// user data based on a specific WebAPI.
type WebAPI interface {
	GetUserID(*entity.User) error
	AddAbsenceStatus(*entity.User, [2]entity.CustomTime) error
}
