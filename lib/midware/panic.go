package midware

import (
	"errors"
	"log"
	"net/http"
)

// Catch ...
func Catch(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var e error
		defer func() {
			r := recover()

			if r != nil {
				switch t := r.(type) {
				case string:
					e = errors.New(t)
				case error:
					e = t
				default:
					e = errors.New("unknown error")
				}

				http.Error(w, e.Error(), http.StatusInternalServerError)

				log.Printf("[panic_handler] %s", e)
			}
		}()

		h(w, r)
	}
}
