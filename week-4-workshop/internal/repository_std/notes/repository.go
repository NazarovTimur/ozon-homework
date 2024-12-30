package notes_repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/model"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		conn: conn,
	}
}

const (
	insertNote = `
		insert into note(title, content, author)
		VALUES ($1, $2, $3)
		returning id`

	insertTag = `
		insert into tag (value)
		values ($1)
		on conflict (value) do nothing`

	insertNoteTag = `
		insert into note_tag(note_id, tag_id)
		select $1, id
		from tag
		where value = $2`

	listNotes = `
		select id, title, content, tags.value tags
		from note n
		         left join lateral (select array_agg(t.value)::text[] value
		                            from tag t
		                                     inner join note_tag nt on t.id = nt.tag_id
		                            where nt.note_id = n.id
		    ) tags on true
		where author = $1
`
)

func (s *Repository) SaveNote(ctx context.Context, note *model.Note) (id int, err error) {
	err = s.conn.QueryRow(ctx, insertNote, note.Title, note.Content, note.Author).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert note: %w", err)
	}
	for _, tag := range note.Tags {
		if _, err = s.conn.Exec(ctx, insertTag, tag); err != nil {
			return 0, fmt.Errorf("insert tag: %w", err)
		}
		if _, err = s.conn.Exec(ctx, insertNoteTag, id, tag); err != nil {
			return 0, fmt.Errorf("insert note tag: %w", err)
		}
	}
	return
}

func (s *Repository) ListNotes(ctx context.Context, author string) (result []*model.Note, err error) {
	rows, err := s.conn.Query(ctx, listNotes, author)
	if err != nil {
		return nil, fmt.Errorf("query list notes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var i model.Note
		if err = rows.Scan(
			&i.Id,
			&i.Title,
			&i.Content,
			&i.Tags,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result = append(result, &i)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return
}
