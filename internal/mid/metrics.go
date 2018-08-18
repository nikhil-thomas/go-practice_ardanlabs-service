package mid

import (
	"context"
	"expvar"
	"net/http"
	"runtime"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"
)

// m contains the global program counters for the application
var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

// Metrics updates program counterss
func Metrics(next web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		v := ctx.Value(web.KeyValues).(*web.Values)

		next(ctx, w, r, params)

		m.req.Add(1)

		if m.req.Value()%100 == 0 {
			m.gr.Set(int64(runtime.NumGoroutine()))
		}

		if v.Error {
			m.err.Add(1)
		}

		return nil
	}

	return h
}
