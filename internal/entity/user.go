package entity

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	DisplayName string
	Email       string
	MobilePhone string
	WorkPhone   string
	ID          uint32
}

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
		validation.Field(
			&u.ID,
			validation.Required),
	)
}
