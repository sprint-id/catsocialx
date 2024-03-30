package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/syarifid/bankx/internal/dto"
	"github.com/syarifid/bankx/internal/ierr"
	"github.com/syarifid/bankx/internal/service"
	response "github.com/syarifid/bankx/pkg/resp"
)

type transactionHandler struct {
	transactionSvc *service.TransactionService
}

func newTransactionHandler(transactionSvc *service.TransactionService) *transactionHandler {
	return &transactionHandler{transactionSvc}
}

func (h *transactionHandler) AddBalance(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqAddBalance

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.transactionSvc.AddBalance(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *transactionHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	balance, err := h.transactionSvc.GetBalance(r.Context(), token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = balance

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}
