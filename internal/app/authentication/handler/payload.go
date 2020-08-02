package handler

import (
	"encoding/json"
	"net/http"
)

type AuthenticationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (payload *AuthenticationRequest) Bind(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(payload); err != nil {
		return err
	}
	return nil
}

type AuthenticationResponse struct {
	AccessToken string `json:"access_token"`
}

func (resp *AuthenticationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

type AuthenticationFailedResponse struct {
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (resp *AuthenticationFailedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(resp.HTTPStatus)
	return nil
}
