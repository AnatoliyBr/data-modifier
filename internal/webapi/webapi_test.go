package webapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/AnatoliyBr/data-modifier/internal/webapi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAppWebAPI_Authenticate(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, _ := token.SignedString([]byte("secret"))

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
		}

		w.Header().Set("Set-Cookie", tokenString)

		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	c := entity.TestCredentials()
	c.AuthURL = s.URL

	webAPI := webapi.NewUserWebAPI(c)

	assert.NoError(t, webAPI.Authenticate())
	assert.Equal(t, tokenString, webAPI.Token)
}
