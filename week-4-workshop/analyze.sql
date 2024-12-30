EXPLAIN(ANALYSE )
SELECT id, title, content, tags.value tags
FROM note n
         LEFT JOIN LATERAL (SELECT ARRAY_AGG(t.value)::text[] value
                            FROM tag t
                                     INNER JOIN note_tag nt ON t.id = nt.tag_id
                            WHERE nt.note_id = n.id
    ) tags ON TRUE
WHERE n.author = 'user';


CREATE INDEX idx_author ON note (author);

EXPLAIN(ANALYSE )
SELECT *
FROM note
WHERE author = 'user1';
