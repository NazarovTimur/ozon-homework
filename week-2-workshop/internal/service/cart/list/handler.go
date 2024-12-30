package list

import (
	"context"
	"gitlab.ozon.dev/14/week-2-workshop/internal/domain"
)

type (
	repository interface {
		ListItem(ctx context.Context, userID int64) []domain.Item
	}

	Handler struct {
		repo repository
	}
)

func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) ListItem(ctx context.Context, userID int64) ([]domain.Item, error) {
	return h.repo.ListItem(ctx, userID), nil
}
