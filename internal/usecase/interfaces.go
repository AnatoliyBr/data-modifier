// Package usecase implements application business logic.
package usecase

import "github.com/AnatoliyBr/data-modifier/internal/entity"

// UseCase is a minimal interface for modifying
// user data based on a third-party system.
type UseCase interface {
	GetUserID(*entity.User) error
	AddAbsenceStatus(*entity.User, [2]entity.CustomTime) error
}
