package repository

import (
	"context"
	"errors"
	"gitlab.ozon.dev/14/week-1-workshop/internal/pkg/reviews/model"
)

type Storage = map[model.Sku][]model.Review

type Repository struct {
	storage Storage
}

func NewReviewRepository(capacity int) *Repository {
	return &Repository{storage: make(Storage, capacity)}
}

func (r *Repository) CreateReview(_ context.Context, review model.Review) (*model.Review, error) {
	if review.SKU < 1 {
		return nil, errors.New("sku must be defined")
	}

	r.storage[review.SKU] = append(r.storage[review.SKU], review)

	return &review, nil
}

func (r *Repository) GetReviews(_ context.Context, sku model.Sku) ([]model.Review, error) {
	if sku < 1 {
		return nil, errors.New("sku must be defined")
	}

	return r.storage[sku], nil
}
