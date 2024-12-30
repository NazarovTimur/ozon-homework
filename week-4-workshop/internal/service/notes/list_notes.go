package notes_usecase

import (
	"context"

	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/model"
)

func (s *Service) ListNotes(ctx context.Context, author string) ([]*model.Note, error) {
	return s.repository.ListNotes(ctx, author)
}
