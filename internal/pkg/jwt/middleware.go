package jwt

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func WrapChiRouterWithAuthorization(r chi.Router) chi.Router {
	tokenAuth := jwtauth.New(Algorithm, Secret, nil)
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)
	return r
}
