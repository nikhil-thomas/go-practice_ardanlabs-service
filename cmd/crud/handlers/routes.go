package handlers

import (
	"net/http"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/mid"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"
)

// API returns a handler for a set of routes
func API(masterDB *db.DB) http.Handler {
	//Create the application
	app := web.New(mid.RequestLogger, mid.ErrorHandler)

	// Bind all the user handlers

	u := User{
		MasterDB: masterDB,
	}

	app.Handle("GET", "/v1/users", u.List)
	app.Handle("POST", "/v1/users", u.Create)
	app.Handle("GET", "/v1/users/:id", u.Retrieve)
	app.Handle("PUT", "/v1/users/:id", u.Update)
	app.Handle("DELETE", "/v1/users:id", u.Delete)
	return app
}
