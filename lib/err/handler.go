package err

import (
	"log"
	"net/http"
	"runtime"
	"strings"
)

type Handler struct {
	ErrorPages map[int]string
	LogFile    string
	*log.Logger
}

func (h Handler) errPage(w http.ResponseWriter, code int) {}

func (h Handler) recovery(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := recover()
		if p == nil {
			return
		}

		var (
			name, file string
			line       int
			pc         [16]uintptr
		)

		// SEE: https://gist.github.com/swdunlop/9629168
		// SEE: https://github.com/mholt/caddy/blob/a5f20829cb87c3474b68e7af4def46e71b8864e0/middleware/errors/errors.go#L24
		n := runtime.Callers(3, pc[:])
		for _, pc := range pc[:n] {
			fn := runtime.FuncForPC(pc)
			if fn == nil {
				continue
			}

			file, line = fn.FileLine(pc)
			name = fn.Name()
			if !strings.HasPrefix(name, "runtime.") {
				break
			}
		}

		h.Printf("[PANIC] %s] %s:%d - %v", r.URL.String(), file, line, p)
		h.errPage(w, http.StatusInternalServerError)
	}
}
