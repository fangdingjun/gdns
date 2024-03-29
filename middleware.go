package main

import (
	"net/http"

	log "github.com/fangdingjun/go-log/v5"
)

type logHandler struct {
	status int
	w      http.ResponseWriter
	size   int
}

func (lh *logHandler) WriteHeader(status int) {
	lh.status = status
	lh.w.WriteHeader(status)
}

func (lh *logHandler) Write(buf []byte) (int, error) {
	lh.size += len(buf)
	return lh.w.Write(buf)
}

func (lh *logHandler) Header() http.Header {
	return lh.w.Header()
}

func (lh *logHandler) Status() int {
	if lh.status != 0 {
		return lh.status
	}
	return 200
}

var _ http.ResponseWriter = &logHandler{}

// LogHandler is a log middleware
// record request log
func LogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				log.Infof("\"%s - %s %s %s\" - %d %d \"%s\"",
					r.RemoteAddr, r.Method, r.RequestURI, r.Proto, 500, 0, r.UserAgent())
			}
		}()

		lh := &logHandler{w: w}
		handler.ServeHTTP(lh, r)
		log.Infof("\"%s - %s %s %s\" - %d %d \"%s\"",
			r.RemoteAddr, r.Method, r.RequestURI, r.Proto, lh.Status(), lh.size, r.UserAgent())
	})
}
