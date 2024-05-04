package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/catsocialx/internal/dto"
	"github.com/sprint-id/catsocialx/internal/ierr"
	"github.com/sprint-id/catsocialx/internal/service"
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

	// Check for missing fields
	for _, field := range expectedFields {
		if _, ok := jsonData[field]; !ok {
			http.Error(w, "missing field in request body: "+field, http.StatusBadRequest)
			return
		}
	}

	// Check for invalid fields
	req.MatchCatId, _ = jsonData["matchCatId"].(string)
	req.UserCatId, _ = jsonData["userCatId"].(string)
	req.Message, _ = jsonData["message"].(string)

	if req.MatchCatId == "" || req.UserCatId == "" || req.Message == "" {
		http.Error(w, "missing field in request body", http.StatusBadRequest)
		return
	}

	if len(req.Message) < 5 || len(req.Message) > 120 {
		http.Error(w, "invalid field in request body: message", http.StatusBadRequest)
		return
	}

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

	w.WriteHeader(http.StatusOK)
}
