package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/pkg/jwt"
)

type TransactionEndpoint struct {
	userSessionHelper domain.UserSessionHelper
}

func NewTransactionEndpoint(userSessionHelper domain.UserSessionHelper) *TransactionEndpoint {
	return &TransactionEndpoint{userSessionHelper: userSessionHelper}
}

func (TransactionEndpoint *TransactionEndpoint) BindRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r = jwt.WrapChiRouterWithAuthorization(r)
		r.Post("/transaction", TransactionEndpoint.HandleCreateTransaction)
		r.Put("/transaction/{id}/verify", TransactionEndpoint.HandleVerifyTransaction)
		r.Get("/transaction/{id}", TransactionEndpoint.HandleGetTransaction)
	})

}

func (transactionEndpoint *TransactionEndpoint) HandleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	request := &CreateTransactionRequest{}
	userSession, err := transactionEndpoint.userSessionHelper.GetFromContext(r.Context())
	if err != nil {
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		})
		return
	}
	log.Println(userSession.ID)
	if err := request.Bind(r); err != nil {
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		})
		return
	}

	render.JSON(w, r, &CreateTransactionSuccess{})
}

func (transactionEndpoint *TransactionEndpoint) HandleVerifyTransaction(w http.ResponseWriter, r *http.Request) {
	request := &VerifyTransactionRequest{}
	id := chi.URLParam(r, "id")
	if err := request.Bind(r); err != nil {
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		})
	}
	if id != "1111" {
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusNotFound,
			Message:  "transaction not found",
		})
		return
	}
	if request.Credential != "123456" {
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  "invalid credential",
		})
		return
	}

	render.JSON(w, r, &VerifyTransactionSuccess{
		TransactionID: id,
	})
}

func (transactionEndpoint *TransactionEndpoint) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("transaction id", id)
	render.JSON(w, r, &GetTransactionSuccess{})
}
