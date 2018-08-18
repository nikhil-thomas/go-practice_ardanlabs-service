package web

import (
	"context"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/pborman/uuid"
)

// TraceIDHeader is the header added to outgoing requests
// this header adds traceID to a request
const TraceIDHeader = "X-Trace-ID"

// Key represents the type of value of the context key
type ctxKey int

// KeyValues is how request stores/retrieves values
const KeyValues ctxKey = 1

// Values represent state of each request
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
	Error      bool
}

// Handler is a type that handle http requests within this app
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error

// App is the entry point into our application
// App configures our context object for each of our http handlers
type App struct {
	*httptreemux.TreeMux
	mw []Middleware
}

// New creates an App value that handle a set of routes for the application
func New(mw ...Middleware) *App {
	return &App{
		TreeMux: httptreemux.New(),
		mw:      mw,
	}
}

// Handle is the mechanism to to mount Handlers for a given HTTP verp and path pair
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware) {

	// function to execute for each request
	handler = wrapMiddleware(wrapMiddleware(handler, mw), a.mw)
	h := func(w http.ResponseWriter, r *http.Request, params map[string]string) {

		// Check the reqeust for an existing TraceId
		// if not generate a new one

		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			traceID = uuid.New()
		}
		v := Values{
			TraceID: traceID,
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)

		// Set traceID on the outgoing requests before anyother header
		// to ensure that the trace id is ALWAYS added to the request regardless of
		w.Header().Set(TraceIDHeader, v.TraceID)

		handler(ctx, w, r, params)
	}
	a.TreeMux.Handle(verb, path, h)
}
