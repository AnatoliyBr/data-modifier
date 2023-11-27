package testserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/gorilla/mux"
)

type TestServer struct {
	employeePath   string
	absencePath    string
	port           string
	basicAuthToken string
	u              *entity.User
	absenceData    *entity.UserAbsenceData
	r              *mux.Router
}

func NewTestServer(employeePath, absencePath, port, token string, u *entity.User, ad *entity.UserAbsenceData) *TestServer {
	s := &TestServer{
		employeePath:   employeePath,
		absencePath:    absencePath,
		port:           port,
		basicAuthToken: token,
		u:              u,
		absenceData:    ad,
		r:              mux.NewRouter(),
	}

	s.configureRouter()

	return s
}

func (s *TestServer) configureRouter() {
	s.r.HandleFunc("/"+s.employeePath, s.handleGetUserID()).Methods(http.MethodPost)
	s.r.HandleFunc("/"+s.absencePath, s.handleAddAbsenceStatus()).Methods(http.MethodPost)
}

func (s *TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *TestServer) StartTestServer() error {
	return http.ListenAndServe(s.port, s)
}

func (s *TestServer) handleGetUserID() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
	}

	type response struct {
		Status string         `json:"status"`
		Data   []*entity.User `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header["Authorization"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, errors.New("incorrect auth header"))
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" || authHeaderParts[1] == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("incorrect auth header"))
			return
		}

		if authHeaderParts[1] != s.basicAuthToken {
			s.error(w, r, http.StatusUnauthorized, errors.New("not authenticate"))
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.Email == s.u.Email {
			resp := &response{
				Status: "OK",
				Data:   []*entity.User{s.u},
			}
			s.respond(w, r, http.StatusOK, resp)
		} else {
			s.error(w, r, http.StatusNotFound, errors.New("user id not found"))
		}
	}
}

func (s *TestServer) handleAddAbsenceStatus() http.HandlerFunc {
	type request struct {
		PersonsIDs []int             `json:"personsIds"`
		DateFrom   entity.CustomTime `json:"dateFrom"`
		DateTo     entity.CustomTime `json:"dateTo"`
	}

	type response struct {
		Status string                    `json:"status"`
		Data   []*entity.UserAbsenceData `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header["Authorization"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, errors.New("incorrect auth header"))
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" || authHeaderParts[1] == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("incorrect auth header"))
			return
		}

		if authHeaderParts[1] != s.basicAuthToken {
			s.error(w, r, http.StatusUnauthorized, errors.New("not authenticate"))
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.PersonsIDs[0] == s.u.ID {
			w.WriteHeader(http.StatusOK)
			resp := &response{
				Status: "OK",
				Data:   []*entity.UserAbsenceData{s.absenceData},
			}
			s.respond(w, r, http.StatusOK, resp)
		} else {
			s.error(w, r, http.StatusNotFound, errors.New("status info not found"))
		}
	}
}

func (s *TestServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *TestServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		enc := json.NewEncoder(w)
		//enc.SetIndent("", "    ")
		enc.Encode(data)
	}
}
