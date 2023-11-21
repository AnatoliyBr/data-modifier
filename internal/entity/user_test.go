package entity_test

import (
	"testing"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *entity.User
		isValid bool
	}{
		{
			name: "valid ru name",
			u: func() *entity.User {
				return entity.TestUser()
			},
			isValid: true,
		},
		{
			name: "valid eng name",
			u: func() *entity.User {
				u := entity.TestUser()
				u.DisplayName = "Ivanov Semen Petrovich"
				return u
			},
			isValid: true,
		},
		{
			name: "valid complex name",
			u: func() *entity.User {
				u := entity.TestUser()
				u.DisplayName = "Иванов-Сидоров Семен Петрович"
				return u
			},
			isValid: true,
		},
		{
			name: "valid mobile phone starts with 8",
			u: func() *entity.User {
				u := entity.TestUser()
				u.MobilePhone = "81234567890"
				return u
			},
			isValid: true,
		},
		{
			name: "invalid name",
			u: func() *entity.User {
				u := entity.TestUser()
				u.DisplayName = "?#@*&%!"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *entity.User {
				u := entity.TestUser()
				u.Email = "invalid"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid mobile phone",
			u: func() *entity.User {
				u := entity.TestUser()
				u.MobilePhone = "invalid"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid work phone",
			u: func() *entity.User {
				u := entity.TestUser()
				u.WorkPhone = "invalid"
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}
