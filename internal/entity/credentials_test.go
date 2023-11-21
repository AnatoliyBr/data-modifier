package entity_test

import (
	"testing"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCredentials_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		c       func() *entity.Credentials
		isValid bool
	}{
		{
			name: "valid",
			c: func() *entity.Credentials {
				return entity.TestCredentials()
			},
			isValid: true,
		},
		{
			name: "invalid ip",
			c: func() *entity.Credentials {
				c := entity.TestCredentials()
				c.IP = "invalid"
				return c
			},
			isValid: false,
		},
		{
			name: "invalid port",
			c: func() *entity.Credentials {
				c := entity.TestCredentials()
				c.Port = "invalid"
				return c
			},
			isValid: false,
		},
		{
			name: "invalid url",
			c: func() *entity.Credentials {
				c := entity.TestCredentials()
				c.AbsenceURL = "invalid"
				return c
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.c().Validate())
			} else {
				assert.Error(t, tc.c().Validate())
			}
		})
	}
}
