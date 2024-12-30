package notes

import (
	"context"

	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/mw"

	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/model"

	servicepb "gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) ListNotes(ctx context.Context, _ *emptypb.Empty) (*servicepb.ListNotesResponse, error) {
	author := mw.GetLogin(ctx)
	ns, err := s.impl.ListNotes(ctx, author)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return repackListToProto(ns), nil
}

func repackListToProto(in []*model.Note) *servicepb.ListNotesResponse {
	result := make([]*servicepb.Note, 0, len(in))
	for _, n := range in {
		result = append(result, &servicepb.Note{
			NoteId: uint64(n.Id),
			Info: &servicepb.NoteInfo{
				Title:   n.Title,
				Content: n.Content,
				Tags:    n.Tags,
			},
		})
	}
	return &servicepb.ListNotesResponse{Notes: result}
}
