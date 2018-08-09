package handlers

import (
	"context"
	"net/http"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/user"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"
	"github.com/pkg/errors"
)

// User represents the User API method handler set
type User struct {
	MasterDB *db.DB
}

// List returns all the existing users in the system
func (u *User) List(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbConn, err := u.MasterDB.Copy()
	if err != nil {
		return errors.Wrapf(web.ErrDBNotConfigured, "")
	}
	defer dbConn.Close()

	cu := user.CreateUser{
		UserType:  1,
		FirstName: "bill",
		LastName:  "kennedy",
		Email:     "bill@ardanlabs.com",
		Company:   "ardan",
	}

	usr, err := user.Create(ctx, dbConn, &cu)
	if err != nil {
		return err
	}

	web.Respond(ctx, w, usr, http.StatusCreated)

	return nil
}
