package notes

import (
	"context"

	desc "gitlab.ozon.dev/go/classroom-14/students/week-3-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ desc.NotesServer = (*Service)(nil)

type Service struct {
	desc.UnimplementedNotesServer
}

func NewService() *Service {
	return &Service{}
}

func (Service) SaveNote(ctx context.Context, request *desc.SaveNoteRequest) (*desc.SaveNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveNote not implemented")
}

func (Service) UpdateNoteByID(ctx context.Context, request *desc.UpdateNoteByIDRequest) (*desc.UpdateNoteByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNoteByID not implemented")
}
