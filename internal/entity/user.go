// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User contains info about user.
type User struct {
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	MobilePhone string `json:"mobilePhone,omitempty"`
	WorkPhone   string `json:"workPhone,omitempty"`
	ID          int    `json:"id,omitempty"`
}

// UserAbsenceData represents user absence status
// info from third-party WebAPI.
type UserAbsenceData struct {
	CreatedDate CustomDate `json:"createdDate,omitempty"`
	DateFrom    CustomTime `json:"dateFrom,omitempty"`
	DateTo      CustomTime `json:"dateTo,omitempty"`
	ID          int        `json:"id,omitempty"`
	PersonID    int        `json:"personId,omitempty"`
	ReasonID    int        `json:"reasonId,omitempty"`
}

// Validate validates fields of the User entity.
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.DisplayName,
			validation.Required,
			validation.By(nameValidation())),
		validation.Field(
			&u.Email,
			validation.Required,
			is.Email),
		validation.Field(
			&u.MobilePhone,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\+\d+$|^\d+$`)),
			validation.Length(10, 12)),
		validation.Field(
			&u.WorkPhone,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\+\d+$|^\d+$`)), validation.Length(1, 12)),
	)
}
