package handler

import (
	"encoding/json"
	"net/http"
)

func (h Handler) Migrate(w http.ResponseWriter, r *http.Request) {
	migrationService, err := NewMigrationService(r.Context(), h.SpotifyAuth, h.AppleClient)
	if err == ErrMissingUserTokens {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req MigrateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	migRes, err := migrationService.Migrate(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respB, err := json.Marshal(migRes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(respB)
}
