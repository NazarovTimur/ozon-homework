package server

import (
	"context"
	"gitlab.ozon.dev/14/week-1-workshop/internal/pkg/reviews/model"
)

type ReviewService interface {
	AddReview(ctx context.Context, review model.Review) (*model.Review, error)
	GetReviews(ctx context.Context, sku model.Sku) ([]model.Review, error)
}

type Server struct {
	reviewService ReviewService
}

func New(reviewService ReviewService) *Server {
	return &Server{reviewService: reviewService}
}
