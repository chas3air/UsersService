-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users(
    id UUID PRIMARY KEY,
    login VARCHAR(50),
    password VARCHAR(50)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS Users;
-- +goose StatementEnd
