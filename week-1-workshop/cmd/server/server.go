package main

import (
	server2 "gitlab.ozon.dev/14/week-1-workshop/internal/app/server"
	"gitlab.ozon.dev/14/week-1-workshop/internal/http/middleware"
	"gitlab.ozon.dev/14/week-1-workshop/internal/pkg/reviews/repository"
	"gitlab.ozon.dev/14/week-1-workshop/internal/pkg/reviews/service"
	"log"
	"net/http"
)

func main() {

	log.Println("app starting")

	reviewRepository := repository.NewReviewRepository(100)
	reviewService := service.NewService(reviewRepository)

	server := server2.New(reviewService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /products/{id}/reviews", server.CreateReview)
	mux.HandleFunc("GET /products/{id}/reviews", server.GetReviews)

	timerMux := middleware.NewTimeMux(mux)
	logMux := middleware.NewLogMux(timerMux)

	log.Println("server starting")

	if err := http.ListenAndServe(":8080", logMux); err != nil {
		panic(err)
	}
}
