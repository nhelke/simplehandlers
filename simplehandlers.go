// A collection of HTTP handler filters
package simplehandlers

import (
	"net/http"
	"strings"
	"net/url"
)


// Function returning a private handler which calls the extracts
// the file extension from the path and adds it to the URL query params
// using the ":extension" key and then call the passed handler's ServeHTTP
// function.
func ExtensionHandler(h http.Handler) extH {
	return extH{h}
}

type extH struct {
	h http.Handler
}

func (h extH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimRight(r.URL.Path, "/")
	trailingSlash := len(r.URL.Path) > len(path)
	if dot := strings.LastIndex(path, "."); dot > 0 && strings.Index(path[dot:], "/") < 0 {
		r.URL.RawQuery = url.Values{":extension": []string{strings.ToLower(path[dot:])}}.Encode() + "&" + r.URL.RawQuery
		r.URL.Path = path[0:dot]
		if trailingSlash {
			r.URL.Path += "/"
		}
	}
	h.h.ServeHTTP(w, r)
}

// HTTP HandlerFunc which returns an Error which is automatically converted
// to a 500 error.
type ErrorHandler func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r) and if a non nil error is returned a 500 error
// is returned to the HTTP client.
func (f ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}