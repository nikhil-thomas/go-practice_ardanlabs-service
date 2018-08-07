package mid

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"
)

// RequestLogger logs inormation on each request
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func RequestLogger(next web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		v := ctx.Value(web.KeyValues).(*web.Values)
		next(ctx, w, r, params)

		log.Printf("%s : (%d) : %s %s -> %s (%s)",
			v.TraceID,
			v.StatusCode,
			r.Method, r.URL.Path,
			r.RemoteAddr, time.Since(v.Now))

		return nil
	}
	return h
}
