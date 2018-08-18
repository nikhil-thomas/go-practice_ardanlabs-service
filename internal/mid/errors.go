package mid

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"
)

// ErrorHandler catches and responds to error
func ErrorHandler(next web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		v := ctx.Value(web.KeyValues).(*web.Values)

		// in the event of panic capture error and send down the stack
		defer func() {
			if r := recover(); r != nil {

				v.Error = true
				//Log the panic
				log.Printf("%s : ERROR : Panic Caught : %s\n", v.TraceID, r)

				//Respond with error
				web.RespondError(ctx, w, errors.New("unhandled"), http.StatusInternalServerError)

				// Print stack
				log.Printf("%s : ERROR : Stacktrace\n%s\n", v.TraceID, debug.Stack())
			}
		}()

		if err := next(ctx, w, r, params); err != nil {

			v.Error = true

			err := errors.Cause(err)

			if err != web.ErrNotFound {
				log.Printf("%s : ERROR : %v\n", v.TraceID, err)
			}

			// respond with error
			web.Error(ctx, w, err)

			// error has been handled hence return no error
			return nil
		}
		return nil
	}

	return h
}
