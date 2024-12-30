package notes_repository

import (
	"context"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"

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

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

func (s *Repository) SaveNote(ctx context.Context, note *model.Note) (id int, err error) {
	sql, args, err := psql.Insert("note").Columns("title", "content", "author").
		Values(note.Title, note.Content, note.Author).
		Suffix("RETURNING \"id\"").
		ToSql()
	err = s.conn.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert note: %w", err)
	}
	for _, tag := range note.Tags {
		sql, args, err = psql.Insert("tag").Columns("value").
			Values(tag).
			Suffix("on conflict (value) do nothing").
			ToSql()
		if _, err = s.conn.Exec(ctx, sql, args...); err != nil {
			return 0, fmt.Errorf("insert tag: %w", err)
		}

		sql, args, err = psql.Insert("note_tag").Columns("note_id", "tag_id").
			Select(psql.Select(strconv.Itoa(id), "id").From("tag").Where(sq.Eq{"value": tag})).
			Values(id, tag).
			ToSql()

		if _, err = s.conn.Exec(ctx, sql, args...); err != nil {
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
