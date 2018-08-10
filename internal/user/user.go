package user

import (
	"context"
	"fmt"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"

	mgo "gopkg.in/mgo.v2"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

const usersCollection = "users"

// Create inserts a new user into the database
func Create(ctx context.Context, dbConn *db.DB, cu *CreateUser) (*User, error) {
	now := time.Now()

	u := User{
		UserID:       bson.NewObjectId().Hex(),
		UserType:     cu.UserType,
		FirstName:    cu.FirstName,
		LastName:     cu.LastName,
		Email:        cu.Email,
		Company:      cu.Company,
		DateCreated:  &now,
		DateModified: &now,
	}

	f := func(collection *mgo.Collection) error {
		return collection.Insert(u)
	}

	if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.users.insert(%s)", db.Query(u)))
	}
	return &u, nil
}

// List retrieves a list of existing uisers from the database
func List(ctx context.Context, dbConn *db.DB) ([]User, error) {
	u := []User{}

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).All(&u)
	}
	if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.users.find()")
	}

	return u, nil
}

// Retrieve gets the specified user from the databse
func Retrieve(ctx context.Context, dbConn *db.DB, userID string) (*User, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.Wrapf(web.ErrInvalidID, "bson:IsObjectIddHex: %s", userID)
	}

	q := bson.M{"user_id": userID}

	var u *User

	f := func(collection *mgo.Collection) error {
		return collection.Find(q).One(&u)
	}

	if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return nil, web.ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.users.find(%s)", db.Query(q)))
	}
	return u, nil
}

// Update replaces a user document in the database
func Update(ctx context.Context, dbConn *db.DB, userID string, cu *CreateUser) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.Wrap(web.ErrInvalidID, "check objectid")
	}

	q := bson.M{"user_id": userID}
	m := bson.M{"$set": cu}

	f := func(collection *mgo.Collection) error {
		return collection.Update(q, m)
	}
	if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return web.ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.users.update(%s, %s)", db.Query(q), db.Query(m)))
	}
	return nil
}

// Delete removes a user from the database
func Delete(ctx context.Context, dbConn *db.DB, userID string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.Wrapf(web.ErrInvalidID, "bson.IsObjectIDHex: %s", userID)
	}

	q := bson.M{"user_id": userID}

	f := func(collection *mgo.Collection) error {
		return collection.Remove(q)
	}
	if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return web.ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.users.remove(%s)", db.Query(q)))
	}
	return nil
}
