-- +goose Up
-- +goose StatementBegin
CREATE TABLE stocks (
    sku BIGINT PRIMARY KEY,
    count INT NOT NULL CHECK (count >= 0)
);

CREATE TABLE reserved (
    sku BIGINT PRIMARY KEY,
    count INT NOT NULL CHECK (count >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stocks;
drop table reserved;
-- +goose StatementEnd
