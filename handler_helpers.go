package main

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/stackus/errors"
	"github.com/stackus/goht"
)

// errorHandler is a middleware that intercepts errors and writes them to the response
//
// The `next` function is expected to return an error if something goes wrong, but
// otherwise it will be the same as a regular http.Handler
func errorHandler(next func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			var stErr errors.HTTPCoder
			if errors.As(err, &stErr) {
				http.Error(w, err.Error(), stErr.HTTPCode())
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// parseRequest is a helper function to parse the query and form values from a request
//
// It returns the query and form values as url.Values
func parseRequest(r *http.Request) (url.Values, url.Values, error) {
	if err := r.ParseForm(); err != nil {
		return nil, nil, errors.ErrBadRequest.Wrap(err, "failed to parse form")
	}

	return r.URL.Query(), r.PostForm, nil
}

func render(ctx context.Context, template goht.Template, w io.Writer) error {
	err := template.Render(ctx, w)
	if err != nil {
		return errors.ErrInternal.Wrap(err, "failed to render template")
	}
	return nil
}
