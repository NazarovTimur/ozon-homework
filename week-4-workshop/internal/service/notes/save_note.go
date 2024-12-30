package notes_usecase

import (
	"context"

	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/model"
)

func (s *Service) SaveNote(ctx context.Context, user string, note *model.Note) (int, error) {
	note.Author = user
	return s.repository.SaveNote(ctx, note)
}
