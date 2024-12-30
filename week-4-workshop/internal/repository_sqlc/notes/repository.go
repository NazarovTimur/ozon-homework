package notes_repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/model"
)

type Repository struct {
	q    Querier
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		q:    New(conn),
		conn: conn,
	}
}

func (s *Repository) SaveNote(ctx context.Context, note *model.Note) (id int, err error) {
	err = pgx.BeginFunc(ctx, s.conn, func(tx pgx.Tx) (err error) {
		repository := New(tx)

		noteID, err := repository.InsertNote(ctx, &InsertNoteParams{
			Title:   note.Title,
			Author:  note.Author,
			Content: &note.Content,
		})
		if err != nil {
			return fmt.Errorf("insert note: %w", err)
		}
		id = int(noteID)
		for _, tag := range note.Tags {
			if err = repository.InsertTag(ctx, tag); err != nil {
				return fmt.Errorf("insert tag: %w", err)
			}
			if err = repository.InsertNoteTag(ctx, &InsertNoteTagParams{
				NoteID: noteID,
				Tag:    tag,
			}); err != nil {
				return fmt.Errorf("insert note tag: %w", err)
			}
		}
		return nil
	})
	return
}

func (s *Repository) ListNotes(ctx context.Context, author string) (result []*model.Note, err error) {
	notes, err := s.q.ListNotes(ctx, author)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	for _, n := range notes {
		note := &model.Note{
			Id:    int(n.ID),
			Title: n.Title,
			Tags:  n.Tags,
		}
		if n.Content != nil {
			note.Content = *n.Content
		}
		result = append(result, note)
	}
	return
}
