package testserver_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/AnatoliyBr/data-modifier/pkg/testserver"
	"github.com/stretchr/testify/assert"
)

func TestTestServer_HandleGetUserID(t *testing.T) {
	type response struct {
		Status string         `json:"status"`
		Data   []*entity.User `json:"data"`
	}

	u := entity.TestUser()

	r := &response{
		Status: "OK",
		Data:   []*entity.User{u},
	}

	s := testserver.NewTestServer(
		"employees",
		"absences",
		":8082",
		"token",
		u,
	)

	testCases := []struct {
		name         string
		payload      interface{}
		authToken    string
		expectedCode int
		isValid      bool
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email": u.Email,
			},
			authToken:    "token",
			expectedCode: http.StatusOK,
			isValid:      true,
		},
		{
			name: "invalid auth header",
			payload: map[string]string{
				"email": u.Email,
			},
			authToken:    "",
			expectedCode: http.StatusBadRequest,
			isValid:      false,
		},
		{
			name: "invalid token",
			payload: map[string]string{
				"email": u.Email,
			},
			authToken:    "invalid",
			expectedCode: http.StatusUnauthorized,
			isValid:      false,
		},
		{
			name:         "invalid payload",
			payload:      "",
			authToken:    "token",
			expectedCode: http.StatusBadRequest,
			isValid:      false,
		},
		{
			name: "user id not found",
			payload: map[string]string{
				"email": "invalid",
			},
			authToken:    "token",
			expectedCode: http.StatusNotFound,
			isValid:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/employees", b)
			req.Header.Add("Authorization", "Basic "+tc.authToken)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

			resp := &response{}
			if tc.isValid {
				json.NewDecoder(rec.Body).Decode(resp)
				assert.Equal(t, r, resp)
			} else {
				json.NewDecoder(rec.Body).Decode(resp)
				assert.NotEqual(t, r, resp)
			}
		})
	}
}
