package middleware

import (
	"log"
	"net/http"
)

type LogMux struct {
	h http.Handler
}

func NewLogMux(h http.Handler) http.Handler {
	return &LogMux{h: h}
}

func (m *LogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("request got")

	m.h.ServeHTTP(w, r)

	log.Printf("request processed")
}
