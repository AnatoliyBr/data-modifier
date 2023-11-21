package webapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
)

type UserWebAPI struct {
	Credentials *entity.Credentials
	Token       string
}

func NewUserWebAPI(c *entity.Credentials) *UserWebAPI {
	return &UserWebAPI{
		Credentials: &entity.Credentials{
			IP:         c.IP,
			Port:       c.Port,
			Login:      c.Login,
			Password:   c.Password,
			AbsenceURL: c.AbsenceURL,
			AuthURL:    c.AuthURL,
		}}
}

func (a *UserWebAPI) Authenticate() error {
	payload := map[string]string{
		"login":    a.Credentials.Login,
		"password": a.Credentials.Password,
	}

	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.Credentials.AuthURL, b)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	a.Token = resp.Header["Set-Cookie"][0]

	if a.Token == "" {
		return errors.New("empty header")
	}

	return nil
}
