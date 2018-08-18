package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

var (
	// ErrNotFound is abstracting the mgo not found error
	ErrNotFound = errors.New("Entity not found")

	// ErrInvalidID occurs when an ID is not in a valid form
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrValidation occurs when there are validtaion errors
	ErrValidation = errors.New("Validation errors occured")

	// ErrNotHealthy occurs when the service is not working properly
	ErrNotHealthy = errors.New("Not Healthy")
)

// JSONError is the response for errors that occur within the API
type JSONError struct {
	Error  string       `json:"error"`
	Fields InvalidError `json:"fields, omitempty"`
}

// Error handles all error responses for the API
func Error(ctx context.Context, w http.ResponseWriter, err error) {
	switch errors.Cause(err) {
	case ErrNotHealthy:
		RespondError(ctx, w, err, http.StatusInternalServerError)
	case ErrNotFound:
		RespondError(ctx, w, err, http.StatusNotFound)
		return
	case ErrValidation, ErrInvalidID:
		RespondError(ctx, w, err, http.StatusBadRequest)
		return
	}

	switch e := errors.Cause(err).(type) {
	case InvalidError:
		v := JSONError{
			Error:  "field validation faliure",
			Fields: e,
		}
		Respond(ctx, w, v, http.StatusBadRequest)
		return
	}
	RespondError(ctx, w, err, http.StatusInternalServerError)
}

// RespondError sends JSON describing the error
func RespondError(ctx context.Context, w http.ResponseWriter, err error, code int) {
	Respond(ctx, w, JSONError{Error: err.Error()}, code)
}

// Respond send json to the client
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, code int) {
	// Set the status code for request logger middleware
	v := ctx.Value(KeyValues).(*Values)
	v.StatusCode = code

	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(code)
		return
	}

	// Marshal data into json string
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Printf("%s : Respond %v Marshalling JSON response", v.TraceID, err)
		RespondError(ctx, w, err, http.StatusInternalServerError)
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Write staus code
	w.WriteHeader(code)

	// send result to client
	w.Write(jsonData)
}
