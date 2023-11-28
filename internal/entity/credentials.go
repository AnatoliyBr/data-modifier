package entity

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Credentials contains info for authentication
// and working with a third-party WebAPI.
type Credentials struct {
	IP          string
	Port        string
	Login       string
	Password    string
	EmployeeURL string
	AbsenceURL  string
}

// Validate validates fields of the Credentials entity.
func (c *Credentials) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.IP, validation.Required, is.IP),
		validation.Field(&c.Port, validation.Required, validation.By(portValidation())),
		validation.Field(&c.Login, validation.Required, validation.Match(regexp.MustCompile(`^[\w]+$`))),
		validation.Field(&c.Password, validation.Required, validation.Length(6, 20)),
		validation.Field(&c.EmployeeURL, validation.Required, is.URL),
		validation.Field(&c.AbsenceURL, validation.Required, is.URL),
	)
}
