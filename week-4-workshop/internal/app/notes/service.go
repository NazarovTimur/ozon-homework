package notes

import (
	"context"

	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/model"
	servicepb "gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/pkg/api/notes/v1"
)

var _ servicepb.NotesServer = (*Service)(nil)

type NoteService interface {
	ListNotes(ctx context.Context, user string) ([]*model.Note, error)
	SaveNote(ctx context.Context, user string, note *model.Note) (int, error)
}

type Service struct {
	servicepb.UnimplementedNotesServer
	impl NoteService
}

func NewService(impl NoteService) *Service {
	return &Service{impl: impl}
}
