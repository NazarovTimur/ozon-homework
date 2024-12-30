-- name: InsertTag :exec
insert into tag (value)
values (@note)
on conflict (value) do nothing;

-- name: InsertNoteTag :exec
insert into note_tag(note_id, tag_id)
select @note_id, id
from tag
where value = @tag;

-- name: InsertNote :one
insert into note(title, content, author)
VALUES (@title, @content, @author)
returning id;

-- name: ListNotes :many
select id, title, content, tags.value tags
from note n
         left join lateral (select array_agg(t.value)::text[] value
                            from tag t
                                     inner join note_tag nt on t.id = nt.tag_id
                            where nt.note_id = n.id
    ) tags on true
where n.author = @author;
