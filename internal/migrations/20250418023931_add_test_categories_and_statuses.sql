-- +goose Up
-- +goose StatementBegin
INSERT INTO categories (name, type) VALUES
    ('Зарплата', 'credit'),
    ('Аванс', 'credit'),
    ('Премия', 'credit'),
    ('Продукты', 'debit'),
    ('Транспорт', 'debit'),
    ('Коммунальные услуги', 'debit'),
    ('Развлечения', 'debit'),
    ('Здоровье', 'debit'),
    ('Одежда', 'debit'),
    ('Прочее', 'debit')
ON CONFLICT (id) DO NOTHING;

INSERT INTO transaction_statuses (name, description) VALUES
    ('Новая', 'Транзакция только что создана'),
    ('В обработке', 'Транзакция находится в процессе обработки'),
    ('Завершена', 'Транзакция успешно завершена'),
    ('Отклонена', 'Транзакция была отклонена'),
    ('Отменена', 'Транзакция была отменена пользователем')
ON CONFLICT (id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM transaction_statuses;
DELETE FROM categories;
-- +goose StatementEnd