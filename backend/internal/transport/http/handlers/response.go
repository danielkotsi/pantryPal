package handlers

import (
	"encoding/json"
	"net/http"

	"pantrypal/backend/internal/transport/http/dto"
)

func ReadJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteAPIError(w http.ResponseWriter, status int, code, message string) {
	var out dto.APIError
	out.Error.Code = code
	out.Error.Message = message
	WriteJSON(w, status, out)
}
