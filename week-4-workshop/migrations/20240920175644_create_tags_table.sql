-- +goose Up
-- +goose StatementBegin
create table tag
(
    id    serial PRIMARY KEY,
    value text not null unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tag;
-- +goose StatementEnd
