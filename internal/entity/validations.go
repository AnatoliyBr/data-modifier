package entity

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func nameValidation() validation.RuleFunc {
	return func(value interface{}) error {
		fullname := strings.ReplaceAll(value.(string), "-", "")
		nameParts := strings.Split(fullname, " ")

		for _, p := range nameParts {
			if err := validation.Validate(p, is.UTFLetter); err != nil {
				return err
			}
		}
		return nil
	}
}

func portValidation() validation.RuleFunc {
	return func(value interface{}) error {
		return validation.Validate(value.(string)[1:], is.Port)
	}
}
