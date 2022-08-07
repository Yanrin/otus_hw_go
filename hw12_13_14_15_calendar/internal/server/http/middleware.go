package httpserver

import (
	"net/http"
	"strings"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		log.Info(strings.Join( // todo: добавить логирование кода ответа
			[]string{
				r.RemoteAddr,
				r.Method,
				r.URL.Path,
				r.Host,
				r.UserAgent(),
			},
			" ",
		))

		next.ServeHTTP(w, r)
	})
}
