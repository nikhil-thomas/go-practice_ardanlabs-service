package handlers

import (
	"context"
	"net/http"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/user"
	"go.opencensus.io/trace"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"
	"github.com/pkg/errors"
)

// check checcks for certain error types and converts them into web errors
func check(err error) error {
	switch errors.Cause(err) {
	case user.ErrNotFound:
		return web.ErrNotFound
	case user.ErrInvalidID:
		return web.ErrInvalidID
	}
	return err
}

// User represents the User API method handler set
type User struct {
	MasterDB *db.DB
}

// List returns all the existing users in the system
func (u *User) List(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctx, span := trace.StartSpan(ctx, "handlers.User.List")
	defer span.End()

	dbConn := u.MasterDB.Copy()
	defer dbConn.Close()

	usrs, err := user.List(ctx, dbConn)

	if err := check(err); err != nil {
		return errors.Wrap(err, "")
	}

	web.Respond(ctx, w, usrs, http.StatusOK)

	return nil
}

// Retrieve returns the specified user from the system
func (u *User) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctx, span := trace.StartSpan(ctx, "handlers.User.retrieve")
	defer span.End()

	dbConn := u.MasterDB.Copy()

	defer dbConn.Close()

	usr, err := user.Retrieve(ctx, dbConn, params["id"])
	if err := check(err); err != nil {
		return errors.Wrapf(err, "Id: %s", params["id"])
	}

	web.Respond(ctx, w, usr, http.StatusOK)
	return nil
}

// Create inserts a new user into the system
func (u *User) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctx, span := trace.StartSpan(ctx, "handlers.User.create")
	defer span.End()

	dbConn := u.MasterDB.Copy()

	defer dbConn.Close()

	var usr user.CreateUser

	if err := web.Unmarshal(r.Body, &usr); err != nil {
		return errors.Wrap(err, "")
	}

	nUsr, err := user.Create(ctx, dbConn, &usr)

	if err := check(err); err != nil {
		return errors.Wrapf(err, "User: %v", &usr)
	}

	web.Respond(ctx, w, nUsr, http.StatusCreated)
	return nil
}

// Update updates the specified user in the system
func (u *User) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctx, span := trace.StartSpan(ctx, "handlers.User.update")
	defer span.End()

	dbConn := u.MasterDB.Copy()

	var usr user.CreateUser

	if err := web.Unmarshal(r.Body, &usr); err != nil {
		return errors.Wrapf(err, "Id: %s User: %+v", params["id"], &usr)
	}

	err := user.Update(ctx, dbConn, params["id"], &usr)
	if err := check(err); err != nil {
		return errors.Wrapf(err, "Id: %s User: %+v", params["id"], &usr)
	}

	web.Respond(ctx, w, nil, http.StatusNoContent)
	return nil
}

// Delete removed the specified user from the system
func (u *User) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ctx, span := trace.StartSpan(ctx, "handlers.User.delete")
	defer span.End()

	dbConn := u.MasterDB.Copy()
	defer dbConn.Close()

	err := user.Delete(ctx, dbConn, params["id"])
	if err := check(err); err != nil {
		return errors.Wrapf(err, "Id: %s", params["id"])
	}
	web.Respond(ctx, w, nil, http.StatusNoContent)
	return nil
}
