// Package testserver implements HTTP server to test datamodifier service.
package testserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/gorilla/mux"
)

// TestServer represents a third-party WebAPI server
// for testing datamodifier server.
type TestServer struct {
	employeePath   string
	absencePath    string
	port           string
	basicAuthToken string
	u              *entity.User
	absenceData    *entity.UserAbsenceData
	r              *mux.Router
}

// NewTestServer returns TestServer with configured routes.
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

// configureRouter registers endpoints.
func (s *TestServer) configureRouter() {
	s.r.HandleFunc("/"+s.employeePath, s.handleGetUserID()).Methods(http.MethodPost)
	s.r.HandleFunc("/"+s.absencePath, s.handleAddAbsenceStatus()).Methods(http.MethodPost)
}

// ServeHTTP handles request.
func (s *TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

// StartTestServer starts listen and handle requests.
func (s *TestServer) StartTestServer() error {
	return http.ListenAndServe(s.port, s)
}

// handleGetUserID handles endpoint for searching user id.
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

// handleAddAbsenceStatus handles endpoint for searching user absence status.
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

// error creates an error response.
func (s *TestServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond creates a response with JSON body.
func (s *TestServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		enc := json.NewEncoder(w)
		//enc.SetIndent("", "    ")
		enc.Encode(data)
	}
}
