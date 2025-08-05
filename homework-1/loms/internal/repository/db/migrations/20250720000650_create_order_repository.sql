-- +goose Up
-- +goose StatementBegin

create table orders
(
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    status text NOT NULL
);

create table items
(
    id bigserial PRIMARY KEY,
    sku bigint NOT NULL UNIQUE,
    count bigint NOT NULL
);

create table orders_items
(
    order_id bigint NOT NULL references orders (id),
    items_id bigint NOT NULL references items (id),
    unique (order_id, items_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table orders CASCADE;
drop table items CASCADE;
drop table orders_items CASCADE;
-- +goose StatementEnd
