-- +goose Up
-- +goose StatementBegin

CREATE TABLE participants (
    part_id SERIAL PRIMARY KEY,
    part_type VARCHAR(10),
    part_name VARCHAR(255),
    part_bank VARCHAR(255),
    part_account VARCHAR(255),
    part_inn VARCHAR(20) UNIQUE,
    part_phone VARCHAR(20)
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    login_name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'user')),
    part_id INTEGER NOT NULL REFERENCES participants(part_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS participants;
-- +goose StatementEnd
