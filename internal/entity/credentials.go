package entity

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Credentials struct {
	IP         string
	Port       string
	Login      string
	Password   string
	AbsenceURL string
	AuthURL    string
}

func (c *Credentials) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.IP, validation.Required, is.IP),
		validation.Field(&c.Port, validation.Required, validation.By(portValidation())),
		validation.Field(&c.Login, validation.Required, validation.Match(regexp.MustCompile(`^[\w]+$`))),
		validation.Field(&c.Password, validation.Required, validation.Length(6, 20)),
		validation.Field(&c.AbsenceURL, validation.Required, is.URL),
		validation.Field(&c.AuthURL, validation.Required, is.URL),
	)
}
