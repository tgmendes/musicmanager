package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/tgmendes/soundfuse/auth"
	"github.com/tgmendes/soundfuse/worker"
)

func (h Handler) Migrate(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewString()
	var req MigrateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userTokens, ok := auth.TokenFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user tokens", http.StatusUnauthorized)
		return
	}

	task := worker.Task{
		ReqID: id,
		F: func(ctx context.Context) error {
			return h.runMigration(ctx, req, userTokens)
		},
	}
	h.Worker.AddTask(&task)
	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) runMigration(ctx context.Context, req MigrateRequest, userTokens auth.CombinedTokens) error {
	migrationService, err := NewMigrationService(ctx, h.SpotifyAuth, h.AppleClient, h.Cache, userTokens)
	if err != nil {
		return err
	}

	migRes, err := migrationService.Migrate(ctx, req)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", migRes)
	return nil
}
