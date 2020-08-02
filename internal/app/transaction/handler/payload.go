package handler

import (
	"encoding/json"
	"math/big"
	"net/http"
)

type TransactionHandlerFailed struct {
	HttpCode int    `json:"-"`
	Message  string `json:"message"`
}

func (resp *TransactionHandlerFailed) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(resp.HttpCode)
	return nil
}

type CreateTransactionRequest struct {
	TransactionCode     string     `json:"transaction_code"`
	Amount              *big.Float `json:"amount"`
	DestinationAccount  string     `json:"destination_account"`
	AuthorizationMethod string     `json:"auth_method"`
}

func (payload *CreateTransactionRequest) Bind(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(payload); err != nil {
		return err
	}
	return nil
}

type CreateTransactionSuccess struct {
	TransactionID string `json:"transaction_id"`
}

func (resp *CreateTransactionSuccess) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusCreated)
	return nil
}

type VerifyTransactionRequest struct {
	Credential string `json:"credential"`
}

func (payload *VerifyTransactionRequest) Bind(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(payload); err != nil {
		return err
	}
	return nil
}

type VerifyTransactionSuccess struct {
	TransactionID string `json:"transaction_id"`
}

func (resp *VerifyTransactionSuccess) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusAccepted)
	return nil
}

type GetTransactionSuccess struct {
	ID                 string     `json:"id"`
	Amount             *big.Float `json:"amount"`
	DestinationAccount string     `json:"destination_account"`
	State              string     `json:"state"`
}

func (resp *GetTransactionSuccess) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	return nil
}
