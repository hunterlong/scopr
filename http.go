package scoper

import "net/http"

type HandlerFunc func(*http.Request) (interface{}, error)

func HandleFunc(handler func(*http.Request) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
