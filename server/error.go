package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// error for bad requests from the caller. Maybe they malformed a url or
// sent a payload with invalid data
func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// error for unexpected server issues like db connections
// this could be because of a bug in our code (like we failed to validate a unique constraint)
// but most likely it's some server specific thing we want to view in the logs
func ErrInternalServerError(err error) render.Renderer {
	log.Err(err).Msg("internal server error")
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error.",
		ErrorText:      "Internal Server Error.", // don't let the client know about internal server errors (i.e. db errors, system level stuff)
	}
}

func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Resource cannot be updated.",
		ErrorText:      err.Error(),
	}
}


var ErrNotFound = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: "Resource not found."}
var ErrForbidden = &ErrResponse{HTTPStatusCode: http.StatusForbidden, StatusText: "Forbidden."}
