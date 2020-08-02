package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tunaiku/mobilebanking/internal/app/authentication"
	"github.com/tunaiku/mobilebanking/internal/app/savings"
	"github.com/tunaiku/mobilebanking/internal/app/transaction"
	"github.com/tunaiku/mobilebanking/internal/app/user"
	"github.com/tunaiku/mobilebanking/internal/pkg/pg"
	"go.uber.org/dig"
)

var (
	container = dig.New()
)

func init() {
	log.Println("register ...")
	transaction.Register(container)
	pg.Register(container)
	authentication.Register(container)
	savings.Register(container)
	user.Register(container)
	container.Provide(func() chi.Router {
		return chi.NewRouter()
	})
}

func invoke() {
	transaction.Invoke(container)
	authentication.Invoke(container)
	savings.Invoke(container)
	user.Invoke(container)
	err := container.Invoke(func(router chi.Router) error {
		log.Println("running server ...")
		return http.ListenAndServe(":8080", router)
	})
	if err != nil {
		log.Fatal(err)
	}
}
