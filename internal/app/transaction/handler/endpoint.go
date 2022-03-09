package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/pkg/jwt"
)

type TransactionEndpoint struct {
	userSessionHelper             domain.UserSessionHelper
	transactionService            domain.MyTransactionService
	accountInformationService     domain.AccountInformationService
	transactionInformationService domain.TransactionInformationService
	otpCredentialManager          domain.OtpCredentialManager
	pinCredentialManger           domain.PinCredentialManager
}

func NewTransactionEndpoint(
	userSessionHelper domain.UserSessionHelper,
	transactionService domain.MyTransactionService,
	accountInformationService domain.AccountInformationService,
	transactionInformationService domain.TransactionInformationService,
	otpCredentialManager domain.OtpCredentialManager,
	pinCredentialManger domain.PinCredentialManager,

) *TransactionEndpoint {
	return &TransactionEndpoint{
		userSessionHelper:             userSessionHelper,
		transactionService:            transactionService,
		accountInformationService:     accountInformationService,
		transactionInformationService: transactionInformationService,
		otpCredentialManager:          otpCredentialManager,
		pinCredentialManger:           pinCredentialManger,
	}
}

func (TransactionEndpoint *TransactionEndpoint) BindRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r = jwt.WrapChiRouterWithAuthorization(r)
		r.Post("/transaction", TransactionEndpoint.HandleCreateTransaction)
		r.Put("/transaction/{id}/verify", TransactionEndpoint.HandleVerifyTransaction)
		r.Get("/transaction/{id}", TransactionEndpoint.HandleGetTransaction)
	})

}

const (
	ERR_INTERNAL_SERVER_MSG                       = "something went wrong"
	ERR_TRANSACTION_CODE_NOT_FOUND_MSG            = "transaction code not found"
	ERR_TRANSACTION_NOT_ALLOWED_MSG               = "transaction not allowed"
	ERR_DESTINATION_ACCOUNT_NOT_FOUND_MSG         = "destination account not found"
	ERR_MINIMUM_AMOUNT_MSG                        = "amount does not reach the minimum transaction amount"
	ERR_AUTHORIZATION_METHOD_NOT_CONFIGURED_MSG   = "authorization method not configured"
	ERR_UNSUPPORTED_AUTHOROZATION_METHOD_MSG      = "unsupported authorization method"
	ERR_TRANSACTION_NOT_FOUND_MSG                 = "transaction not found"
	ERR_VERIFICATION_PROCESS_ALREADY_HAPPENED_MSG = "verification process already happened"
	ERR_INVALID_CREDENTIAL_MSG                    = "invalid credential"
)

// HandleCreateTransaction
// Get the User Session ( if the user is valid user )
// Decode the request body
// Get the transaction Detail by the TransactionCode (if valid then proceed)
// Check if the Source User is able for the Transaction
// Check if the Destination User is able for the Transaction
// Validate the minimum transaction amount // if min amt = 2000 then transaction of 2000 should be permitted // one test case failed because of this
// Check if the authorization method is configured for the source user (pin or otp in this case)
// Create the transaction in the database with flag domain.WaitAuthorization
// Send the Otp/Pin to the source user (printing on console in this case)
// Send the response to the user with http 201 created
func (transactionEndpoint *TransactionEndpoint) HandleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	request := &CreateTransactionRequest{}
	userSession, err := transactionEndpoint.userSessionHelper.GetFromContext(r.Context())
	if err != nil {
		log.Println("userSession Error:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}
	if err := request.Bind(r); err != nil {
		log.Println("requestDecoding Error:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}

	// get the transaction details
	transactionDetail, err := transactionEndpoint.transactionInformationService.FindTransactionDetailByCode(request.TransactionCode)
	if err != nil {
		log.Println("find transaction deatil Error:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  ERR_TRANSACTION_CODE_NOT_FOUND_MSG,
		})
		return
	}

	// check if the transaction is allowed for user
	isTransactionValidForUser, err := transactionEndpoint.IfTransactionIsAllowedForUser(userSession.AccountReference, request.TransactionCode)
	if err != nil {
		log.Println("source user valid transaction Error:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}
	if !isTransactionValidForUser {
		log.Println("source user invalid for transaction:")
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  ERR_TRANSACTION_NOT_ALLOWED_MSG,
		})
		return
	}

	// check if the destination account is valid and valid for trasaction
	isDestinationAccExist := transactionEndpoint.accountInformationService.IsAccountExists(request.DestinationAccount)
	if !isDestinationAccExist {
		log.Println("destination account doesnot exist")
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  ERR_DESTINATION_ACCOUNT_NOT_FOUND_MSG,
		})
		return
	}
	isTransactionValidForUser = false
	isTransactionValidForUser, err = transactionEndpoint.IfTransactionIsAllowedForUser(request.DestinationAccount, request.TransactionCode)
	if err != nil {
		log.Println("destination user valid for transaction err:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}
	if !isTransactionValidForUser {
		log.Println("destination user invalid for transaction")
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  ERR_TRANSACTION_NOT_ALLOWED_MSG,
		})
		return
	}

	// validate the minimum transaction amount
	if request.Amount.Cmp(transactionDetail.MinimumAmount) == -1 { // if min amt = 2000 then transaction of 2000 should be permitted // one test case failed because of this
		log.Println("transaction amount is less than minimum transactable amount")
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  ERR_MINIMUM_AMOUNT_MSG,
		})
		return
	}

	// check the authorization method configured for user
	var authMethod domain.AuthorizationMethod
	switch request.AuthorizationMethod {
	case "otp":
		if !userSession.User.ConfiguredTransactionCredential.IsOtpConfigured() {
			log.Println("otp credential not configured")
			render.Render(w, r, &TransactionHandlerFailed{
				HttpCode: http.StatusBadRequest,
				Message:  ERR_AUTHORIZATION_METHOD_NOT_CONFIGURED_MSG,
			})
			return
		}
		authMethod = domain.OtpAuthorization
		// transactionEndpoint.otpCredentialManager.RequestNewOtp(userSession.ID)
	case "pin":
		if !userSession.User.ConfiguredTransactionCredential.IsPinConfigured() {
			log.Println("pin credential not configured")
			render.Render(w, r, &TransactionHandlerFailed{
				HttpCode: http.StatusBadRequest,
				Message:  ERR_AUTHORIZATION_METHOD_NOT_CONFIGURED_MSG,
			})
			return
		}
		authMethod = domain.PinAuthorization
		// print the pin
		// log.Println("pin is", userSession.ConfiguredTransactionCredential.Pin.Pin)
	default:
		log.Println(ERR_UNSUPPORTED_AUTHOROZATION_METHOD_MSG)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  ERR_UNSUPPORTED_AUTHOROZATION_METHOD_MSG,
		})
		return
	}

	amount, _ := request.Amount.Float64() // accuracy is ignored for the assignment it should be considered in real scenerios
	tx := &domain.Transaction{
		ID:                  uuid.New().String(), // random uuid
		UserID:              userSession.ID,
		State:               domain.WaitAuthorization,
		AuthorizationMethod: authMethod,
		TransactionCode:     request.TransactionCode,
		Amount:              amount,
		SourceAccount:       userSession.AccountReference,
		DestinationAccount:  request.DestinationAccount,
		CreatedAt:           time.Now(),
	}

	// call the transaction service layer to create the transaction
	transactionID, err := transactionEndpoint.transactionService.CreateTransaction(tx)
	if err != nil {
		log.Println("create transaction error:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}

	// generate and send the otp to user in actual case
	switch request.AuthorizationMethod {
	case "otp":
		transactionEndpoint.otpCredentialManager.RequestNewOtp(userSession.ID)
	case "pin":
		// print the pin
		log.Println("pin is", userSession.ConfiguredTransactionCredential.Pin.Pin)
	}

	// if the otp generation or sending failed then delete the transaction or mark as failed
	// this scenerio will not occur in this case

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, &CreateTransactionSuccess{TransactionID: transactionID})
}

func (transactionEndpoint *TransactionEndpoint) HandleVerifyTransaction(w http.ResponseWriter, r *http.Request) {
	request := &VerifyTransactionRequest{}
	id := chi.URLParam(r, "id")
	if err := request.Bind(r); err != nil {
		log.Println("request decodeing err:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}

	// get the user session
	userSession, err := transactionEndpoint.userSessionHelper.GetFromContext(r.Context())
	if err != nil {
		log.Println("userSession Error:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}

	// get the transaction details
	tx, err := transactionEndpoint.transactionService.ReadTransaction(id)
	if err != nil {
		log.Println("get transaction error Error:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}
	// check if the transaction not found
	if tx == nil {
		log.Println(ERR_TRANSACTION_NOT_FOUND_MSG)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusNotFound,
			Message:  ERR_TRANSACTION_NOT_FOUND_MSG,
		})
		return
	}

	// check if the transaction is owned by the requesting user
	if userSession.AccountReference != tx.SourceAccount {
		log.Println("transaction doesnot belongs to user")
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError, // should be http.StatusForbidden
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}

	// check if the transaction verification already done
	if tx.State != domain.WaitAuthorization && tx.State != domain.UnknownTransactionStatus {
		log.Println("transaction verification already done")
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  ERR_VERIFICATION_PROCESS_ALREADY_HAPPENED_MSG,
		})
		return
	}

	// check the auth method pin or otp
	switch tx.AuthorizationMethod {
	case domain.OtpAuthorization:
		err = transactionEndpoint.otpCredentialManager.Validate(userSession.ID, request.Credential)
	case domain.PinAuthorization:
		err = transactionEndpoint.pinCredentialManger.Validate(userSession.ID, request.Credential)
	default: // this default case will never be executed but written only for roboustness
		log.Println("transaction authorization method not supported")
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  ERR_UNSUPPORTED_AUTHOROZATION_METHOD_MSG,
		})
		return
	}
	if err != nil {
		// set the transaction status to failed
		tx.State = domain.Failed
		err = transactionEndpoint.transactionService.UpdateTransaction(tx)
		if err != nil {
			log.Println("transaction state update to failed err:", err)
			render.Render(w, r, &TransactionHandlerFailed{
				HttpCode: http.StatusInternalServerError,
				Message:  ERR_INTERNAL_SERVER_MSG,
			})
			return
		}
		log.Println("credential validation failed:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusBadRequest,
			Message:  ERR_INVALID_CREDENTIAL_MSG,
		})
		return
	}

	// update the transaction status to success
	tx.State = domain.Success
	err = transactionEndpoint.transactionService.UpdateTransaction(tx)
	if err != nil {
		log.Println("transaction status update to success err:", err)
		render.Render(w, r, &TransactionHandlerFailed{
			HttpCode: http.StatusInternalServerError,
			Message:  ERR_INTERNAL_SERVER_MSG,
		})
		return
	}

	// call the core banking api or any pub sub mechanism
	// then return the response

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	render.JSON(w, r, &VerifyTransactionSuccess{
		TransactionID: id,
	})
}

func (transactionEndpoint *TransactionEndpoint) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("transaction id", id)
	render.JSON(w, r, &GetTransactionSuccess{})
}
