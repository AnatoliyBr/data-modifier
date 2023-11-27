package webapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
)

type UserWebAPI struct {
	Credentials        *entity.Credentials
	AbsencesReasonList map[int]*entity.Reason
	BasicAuthToken     string
}

func NewUserWebAPI(c *entity.Credentials) *UserWebAPI {
	a := &UserWebAPI{
		Credentials: &entity.Credentials{
			IP:          c.IP,
			Port:        c.Port,
			Login:       c.Login,
			Password:    c.Password,
			EmployeeURL: c.EmployeeURL,
			AbsenceURL:  c.AbsenceURL,
		},
		AbsencesReasonList: entity.ReasonList}

	a.basicAuth()

	return a
}

func (a *UserWebAPI) basicAuth() {
	a.BasicAuthToken = base64.StdEncoding.EncodeToString(
		[]byte(a.Credentials.Login + ":" + a.Credentials.Password))
}

func (a *UserWebAPI) GetUserID(u *entity.User) error {
	type request struct {
		Email string `json:"email"`
	}

	payload := &request{
		Email: u.Email,
	}

	type response struct {
		Status string         `json:"status"`
		Data   []*entity.User `json:"data"`
	}

	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.Credentials.EmployeeURL, b)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+a.BasicAuthToken)

	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Do(req)
	if err != nil {
		return err
	}

	resp := &response{}
	if err := json.NewDecoder(r.Body).Decode(resp); err != nil {
		return err
	}

	if resp.Status != "OK" {
		return errors.New("not found")
	}

	u.ID = resp.Data[0].ID

	return nil
}

func (a *UserWebAPI) AddAbsenceStatus(u *entity.User, p [2]entity.CustomTime) error {
	type request struct {
		DateFrom   entity.CustomTime `json:"dateFrom"`
		DateTo     entity.CustomTime `json:"dateTo"`
		PersonsIDs []int             `json:"personsIds"`
	}

	payload := &request{
		PersonsIDs: []int{u.ID},
		DateFrom:   p[0],
		DateTo:     p[1],
	}

	type response struct {
		Status string                    `json:"status"`
		Data   []*entity.UserAbsenceData `json:"data"`
	}

	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.Credentials.AbsenceURL, b)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+a.BasicAuthToken)

	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Do(req)
	if err != nil {
		return err
	}

	resp := &response{}
	if err := json.NewDecoder(r.Body).Decode(resp); err != nil {
		return err
	}

	if resp.Status != "OK" {
		return errors.New("not found")
	}

	e := "❓"

	if _, ok := a.AbsencesReasonList[resp.Data[0].ReasonID]; ok {
		e = a.AbsencesReasonList[resp.Data[0].ReasonID].Emoji
	}

	u.DisplayName = u.DisplayName + " " + e

	return nil
}
