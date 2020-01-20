package scopr

import "net/http"

type Handler func(http.ResponseWriter, *http.Request) (interface{}, string)

func HandlerFunc(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, scope := handler(w, r)
		sendData := New(data, scope)
		NewEncoder(w).Write(sendData)
	})
}
