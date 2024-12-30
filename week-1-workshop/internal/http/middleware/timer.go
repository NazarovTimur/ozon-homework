package middleware

import (
	"log"
	"net/http"
	"time"
)

func Timer(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(now time.Time) {
			log.Printf("handler spent %s", time.Since(now))
		}(time.Now())

		h(w, r)
	}
}

type TimerMux struct {
	h http.Handler
}

func NewTimeMux(h http.Handler) http.Handler {
	return &TimerMux{h: h}
}

func (m *TimerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(now time.Time) {
		log.Printf("handler spent %s", time.Since(now))
	}(time.Now())

	m.h.ServeHTTP(w, r)
}
