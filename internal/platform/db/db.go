package db

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

// ErrInvalidDBProvided is returned in the event that an uninitialized db
// is used to perform actions against
var ErrInvalidDBProvided = errors.New("invalid DB provided")

// DB is collection of support for different DB technologies
type DB struct {
	// MongoDB Support
	database *mgo.Database
	session  *mgo.Session
}

// New returns a new DB value for use with MongoDB based on a registered master session
func New(url string, timeout time.Duration) (*DB, error) {
	//set default timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	// create a session which maintains a pool of socket connections
	// to our MongoDB
	ses, err := mgo.DialWithTimeout(url, timeout)
	if err != nil {
		return nil, errors.Wrapf(err, "mgo.DialWithTimeout: %s,%v", url, timeout)
	}

	ses.SetMode(mgo.Monotonic, true)

	db := DB{
		database: ses.DB(""),
		session:  ses,
	}

	return &db, nil
}

// Close closes the DB value being used with MongoDB
func (db *DB) Close() {
	db.session.Close()
}

// Copy returns a new DB value for use with MongoDB based on master session
func (db *DB) Copy() *DB {
	ses := db.session.Copy()
	newDB := DB{
		database: ses.DB(""),
		session:  ses,
	}

	return &newDB
}

// Execute is used to execute MongoDB Commands
func (db *DB) Execute(collName string, f func(*mgo.Collection) error) error {
	if db == nil || db.session == nil {
		return errors.Wrap(ErrInvalidDBProvided, "db == nil || db.session == nil")
	}
	return f(db.database.C(collName))
}

// ExecuteTimeout is used to execute MongoDB commands witha a timeout
func (db *DB) ExecuteTimeout(timeout time.Duration, collName string, f func(*mgo.Collection) error) error {
	if db == nil || db.session == nil {
		return errors.Wrap(ErrInvalidDBProvided, "db == nil || db.session == nil")
	}
	db.session.SetSocketTimeout(timeout)
	return f(db.database.C(collName))
}

// StatusCheck validates the SDB status is good
func (db *DB) StatusCheck() error {
	return nil
}

// Query provides a string version of the value
func Query(value interface{}) string {
	json, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(json)
}
