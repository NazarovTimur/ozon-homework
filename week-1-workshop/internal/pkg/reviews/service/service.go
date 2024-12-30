package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.ozon.dev/14/week-1-workshop/internal/pkg/reviews/model"
)

type ReviewRepository interface {
	CreateReview(_ context.Context, review model.Review) (*model.Review, error)
	GetReviews(_ context.Context, sku model.Sku) ([]model.Review, error)
}

type ReviewService struct {
	repository ReviewRepository
}

func NewService(repository ReviewRepository) *ReviewService {
	return &ReviewService{repository: repository}
}

func (s *ReviewService) AddReview(ctx context.Context, review model.Review) (*model.Review, error) {
	if review.SKU < 1 || len(review.Comment) == 0 || review.UserID == uuid.Nil {
		return nil, errors.New("fail validation")
	}

	return s.repository.CreateReview(ctx, review)
}

func (s *ReviewService) GetReviews(ctx context.Context, sku model.Sku) ([]model.Review, error) {
	if sku < 1 {
		return nil, errors.New("fail validation")
	}

	return s.repository.GetReviews(ctx, sku)
}
