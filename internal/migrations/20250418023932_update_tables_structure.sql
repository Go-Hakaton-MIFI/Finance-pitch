-- +goose Up
-- +goose StatementBegin
-- Добавляем колонку type в таблицу categories
ALTER TABLE categories ADD COLUMN IF NOT EXISTS type VARCHAR(50) CHECK (type IN ('credit', 'debit'));

-- Добавляем колонку description в таблицу transaction_statuses
ALTER TABLE transaction_statuses ADD COLUMN IF NOT EXISTS description TEXT;

-- Обновляем существующие записи в categories
UPDATE categories SET type = 'debit' WHERE name IN ('Продукты', 'Транспорт', 'Коммунальные услуги', 'Развлечения', 'Здоровье', 'Одежда', 'Прочее');
UPDATE categories SET type = 'credit' WHERE name IN ('Зарплата', 'Аванс', 'Премия');

-- Обновляем существующие записи в transaction_statuses
UPDATE transaction_statuses SET description = 'Транзакция только что создана' WHERE name = 'Новая';
UPDATE transaction_statuses SET description = 'Транзакция находится в процессе обработки' WHERE name = 'В обработке';
UPDATE transaction_statuses SET description = 'Транзакция успешно завершена' WHERE name = 'Выполнена';
UPDATE transaction_statuses SET description = 'Транзакция была отменена пользователем' WHERE name = 'Отменена';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories DROP COLUMN IF EXISTS type;
ALTER TABLE transaction_statuses DROP COLUMN IF EXISTS description;
-- +goose StatementEnd 