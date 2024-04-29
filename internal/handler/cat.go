package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/catsocialx/internal/dto"
	"github.com/sprint-id/catsocialx/internal/ierr"
	"github.com/sprint-id/catsocialx/internal/service"
	response "github.com/sprint-id/catsocialx/pkg/resp"
)

type catHandler struct {
	catSvc *service.CatService
}

func newCatHandler(catSvc *service.CatService) *catHandler {
	return &catHandler{catSvc}
}

func (h *catHandler) AddCat(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqAddCat
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// {
	// 	"name": "", // not null, minLength 1, maxLength 30
	// 	"race": "", /** not null, enum of:
	// 			- "Persian"
	// 			- "Maine Coon"
	// 			- "Siamese"
	// 			- "Ragdoll"
	// 			- "Bengal"
	// 			- "Sphynx"
	// 			- "British Shorthair"
	// 			- "Abyssinian"
	// 			- "Scottish Fold"
	// 			- "Birman" */
	// 	"sex": "", // not null, enum of: "male" / "female"
	// 	"ageInMonth": 1, // not null, min: 1, max: 120082
	// 	"description":"" // not null, minLength 1, maxLength 200
	// 	"imageUrls":[ // not null, minItems: 1, items: not null, should be url
	// 		"","",""
	// 	]
	// }

	// Check for unexpected fields
	expectedFields := []string{"name", "race", "sex", "ageInMonth", "description", "imageUrls"}
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

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.catSvc.AddCat(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *catHandler) GetCat(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetCat

	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	cats, err := h.catSvc.GetCat(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = cats

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}

// The contains function checks if a slice contains a string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
