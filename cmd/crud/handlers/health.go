package handlers

import (
	"context"
	"net/http"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"
)

// Health represents the User API method handler set
type Health struct {
	MasterDB *db.DB
}

// Check validates the service is ready and healthy to accept requests
func (h *Health) Check(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	dbConn := h.MasterDB.Copy()
	defer dbConn.Close()

	if err := dbConn.StatusCheck(); err != nil {
		return err
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	web.Respond(ctx, w, data, http.StatusOK)
	return nil
}
