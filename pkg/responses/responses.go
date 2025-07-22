package responses

import (
	"goAuthExample/pkg/database"
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // raw error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status_text"` // user-level status message
	ErrorText  string `json:"error_text"`  // application-level error message
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

type PostResponse struct {
	Message        string `json:"message"`
	HTTPStatusCode int    `json:"-"`
}

type getRenderer[T any] struct {
	Data           T   `json:"data"`
	HTTPStatusCode int `json:"-"`
}

func (rep getRenderer[T]) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, rep.HTTPStatusCode)
	return nil
}
func (c *PostResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, c.HTTPStatusCode)
	return nil
}

func PostResponseRender(message string) render.Renderer {
	return &PostResponse{
		Message:        message,
		HTTPStatusCode: 200,
	}
}

func UserResponseRender(user *database.User) render.Renderer {
	return &getRenderer[*database.User]{
		Data:           user,
		HTTPStatusCode: 200,
	}
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Server error",
		ErrorText:      err.Error(),
	}
}
