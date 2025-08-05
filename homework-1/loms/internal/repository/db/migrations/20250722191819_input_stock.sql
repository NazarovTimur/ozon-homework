-- +goose Up
-- +goose StatementBegin
INSERT INTO reserved (sku, count) VALUES (773297411, 10),
                                         (1002, 20),
                                         (1003, 30),
                                         (1004, 40),
                                         (1005, 50)
    ON CONFLICT (sku) DO NOTHING;

INSERT INTO stocks (sku, count) VALUES (773297411, 150),
                                       (1002, 200),
                                       (1003, 250),
                                       (1004, 30),
                                       (1005, 350)
    ON CONFLICT (sku) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM reserved WHERE sku IN (773297411, 1002, 1003, 1004, 1005);
DELETE FROM stocks WHERE sku IN (773297411, 1002, 1003, 1004, 1005);
-- +goose StatementEnd
