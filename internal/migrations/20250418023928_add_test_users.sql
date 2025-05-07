-- +goose Up
-- +goose StatementBegin

-- Добавляем тестового участника
INSERT INTO participants (part_type, part_name, part_bank, part_account, part_inn, part_phone)
VALUES ('INDIVIDUAL', 'Тестовый Пользователь', 'Сбербанк', '40817810099910004312', '123456789012', '+79001234567');

-- Добавляем тестового пользователя
INSERT INTO users (login_name, password, role, part_id)
VALUES ('test_user', '$2a$10$abcdefghijklmnopqrstuvwxyz123456', 'user', 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE login_name = 'test_user';
DELETE FROM participants WHERE part_name = 'Тестовый Пользователь';
-- +goose StatementEnd 