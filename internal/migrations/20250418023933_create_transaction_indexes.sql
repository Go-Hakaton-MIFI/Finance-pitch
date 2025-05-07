-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_transactions_user_type ON transactions(user_type);
CREATE INDEX idx_transactions_category_id ON transactions(category_id);
CREATE INDEX idx_transactions_status_id ON transactions(status_id);
CREATE INDEX idx_transactions_date_time ON transactions(date_time);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_transactions_user_type;
DROP INDEX IF EXISTS idx_transactions_category_id;
DROP INDEX IF EXISTS idx_transactions_status_id;
DROP INDEX IF EXISTS idx_transactions_date_time;
DROP INDEX IF EXISTS idx_transactions_created_at;
-- +goose StatementEnd 