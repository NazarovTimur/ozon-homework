-- +goose Up
-- +goose StatementBegin
create table note
(
    id      serial PRIMARY KEY,
    title   text not null,
    author  text not null,
    content text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table note;
-- +goose StatementEnd
