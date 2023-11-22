package webapi_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/AnatoliyBr/data-modifier/internal/webapi"
	"github.com/stretchr/testify/assert"
)

func TestAppWebAPI_GetUserID(t *testing.T) {
	u1 := entity.TestUser()
	u2 := entity.TestUser()
	u2.ID = 0

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Email string `json:"email"`
		}

		authHeader, ok := r.Header["Authorization"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors.New("incorrect auth header"))
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" || authHeaderParts[1] == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors.New("incorrect auth header"))
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		if req.Email == u1.Email {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "OK",
				"data":   []*entity.User{u1},
			})
		}
	}))
	defer s.Close()

	c := entity.TestCredentials()
	c.EmployeeURL = s.URL

	webAPI := webapi.NewUserWebAPI(c)

	assert.NoError(t, webAPI.GetUserID(u2))
	assert.Equal(t, u1.ID, u2.ID)
}

func TestAppWebAPI_AddAbsenceStatus(t *testing.T) {
	u1 := entity.TestUser()
	u2 := entity.TestUser()
	absenceData := entity.TestUserAbsenceData()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			PersonsIDs []int             `json:"personsIds"`
			DateFrom   entity.CustomTime `json:"dateFrom"`
			DateTo     entity.CustomTime `json:"dateTo"`
		}

		type response struct {
			Status string                    `json:"status"`
			Data   []*entity.UserAbsenceData `json:"data"`
		}

		authHeader, ok := r.Header["Authorization"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors.New("incorrect auth header"))
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" || authHeaderParts[1] == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors.New("incorrect auth header"))
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		if req.PersonsIDs[0] == u1.ID {
			w.WriteHeader(http.StatusOK)
			resp := &response{
				Status: "OK",
				Data:   []*entity.UserAbsenceData{absenceData},
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer s.Close()

	c := entity.TestCredentials()
	c.AbsenceURL = s.URL

	webAPI := webapi.NewUserWebAPI(c)

	assert.NoError(t, webAPI.AddAbsenceStatus(
		u1,
		[2]entity.CustomTime{
			absenceData.DateFrom,
			absenceData.DateTo}))
	assert.NotEqual(t, u1.DisplayName, u2.DisplayName)
	assert.Equal(t, u1.DisplayName, u2.DisplayName+" "+entity.ReasonList[absenceData.ReasonID].Emoji)
}
