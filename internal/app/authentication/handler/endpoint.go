package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/micro/go-micro/v3/errors"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
)

type AuthenticationEndpoint struct {
	authenticationService domain.AuthenticationService
}

func NewAuthenticationEndpoint(authenticationService domain.AuthenticationService) *AuthenticationEndpoint {
	return &AuthenticationEndpoint{authenticationService: authenticationService}
}

func (endpoint AuthenticationEndpoint) BindRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/authenticate", endpoint.HandleAuthenticationFlow)
	})
}

func (endpoint AuthenticationEndpoint) HandleAuthenticationFlow(w http.ResponseWriter, r *http.Request) {
	request := new(AuthenticationRequest)
	if err := request.Bind(r); err != nil {
		render.JSON(w, r, &AuthenticationFailedResponse{Message: err.Error(), HTTPStatus: http.StatusInternalServerError})
		return
	}
	result, err := endpoint.authenticationService.Authenticate(request.Username, request.Password)
	if err != nil {
		switch v := err.(type) {
		case *errors.Error:
			log.Print(v.Error())
			render.Render(w, r, &AuthenticationFailedResponse{Message: v.Detail, HTTPStatus: int(v.Code)})
			return
		default:
			log.Print(v.Error())
			render.Render(w, r, &AuthenticationFailedResponse{Message: v.Error(), HTTPStatus: http.StatusInternalServerError})
			return
		}
	}
	render.JSON(w, r, &AuthenticationResponse{AccessToken: result.AccessToken})
}
