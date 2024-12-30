-- +goose Up
-- +goose StatementBegin
create table note_tag
(
    note_id int not null references note (id),
    tag_id  int not null references tag (id),
    unique (note_id, tag_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table note_tag;
-- +goose StatementEnd
