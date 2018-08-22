package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"

	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

var idkey tag.Key

func init() {
	var err error
	if idkey, err = tag.NewKey("idKey"); err != nil {
		log.Fatal(err)
	}
}

// Trace middleware updates spans
func Trace(next web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		v := ctx.Value(web.KeyValues).(*web.Values)

		ctx, err := tag.New(ctx,
			tag.Insert(idkey, "testing tag"),
		)

		if err != nil {
			log.Println("midware : ERROR :", err)
		}

		// Add a SPAN for this request
		ctx, span := trace.StartSpan(ctx, v.TraceID)
		defer span.End()

		next(ctx, w, r, params)

		return nil
	}

	return h
}
