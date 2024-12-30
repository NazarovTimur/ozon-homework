package notes

import (
	"context"

	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/mw"

	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/model"

	servicepb "gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) SaveNote(ctx context.Context, in *servicepb.SaveNoteRequest) (*servicepb.SaveNoteResponse, error) {
	author := mw.GetLogin(ctx)
	id, err := s.impl.SaveNote(ctx, author, repackNote(in))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &servicepb.SaveNoteResponse{NoteId: uint64(id)}, nil
}

func repackNote(in *servicepb.SaveNoteRequest) *model.Note {
	return &model.Note{
		Title:   in.Info.Title,
		Content: in.Info.Content,
		Tags:    in.Info.Tags,
	}
}
