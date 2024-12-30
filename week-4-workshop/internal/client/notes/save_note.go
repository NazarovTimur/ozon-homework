package notes

import (
	"context"

	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/model"
	desc "gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc/metadata"
)

func (c *Client) SaveNote(ctx context.Context, note *model.Note) (int, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "x-auth", "user")
	response, err := c.client.SaveNote(ctx, &desc.SaveNoteRequest{Info: &desc.NoteInfo{
		Title:   note.Title,
		Content: note.Content,
	}})
	if err != nil {
		return 0, err
	}
	return int(response.NoteId), err
}
