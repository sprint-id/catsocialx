package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/catsocialx/internal/dto"
	"github.com/sprint-id/catsocialx/internal/ierr"
	"github.com/sprint-id/catsocialx/internal/service"
	response "github.com/sprint-id/catsocialx/pkg/resp"
)

type matchHandler struct {
	matchSvc *service.MatchService
}

func newMatchHandler(matchSvc *service.MatchService) *matchHandler {
	return &matchHandler{matchSvc}
}

// ReqMatchCat is a struct to represent request payload for match cat
// {
// 	"matchCatId": "",
// 	"userCatId": "",
// 	"message": "" // not null, minLength: 5, maxLength: 120
// }

func (h *matchHandler) MatchCat(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqMatchCat
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"matchCatId", "userCatId", "message"}
	for key := range jsonData {
		if !contains(expectedFields, key) {
			http.Error(w, "unexpected field in request body: "+key, http.StatusBadRequest)
			return
		}
	}

	// Convert the jsonData map into the req struct
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// show request
	fmt.Printf("MatchCat request: %+v\n", req)

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.matchSvc.MatchCat(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *matchHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	matches, err := h.matchSvc.GetMatch(r.Context(), token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	// show response
	// fmt.Printf("GetMatch response: %+v\n", matches)

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = matches

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}
