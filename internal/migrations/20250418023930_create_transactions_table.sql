-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    CREATE TYPE transaction_type AS ENUM ('income', 'expense');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('credit', 'debit')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transaction_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_type VARCHAR(50) NOT NULL,
    date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    trans_type VARCHAR(50) NOT NULL,
    amount DECIMAL(15,5) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    status_id INTEGER REFERENCES transaction_statuses(id),
    sender_bank VARCHAR(255),
    receiver_inn VARCHAR(12),
    receiver_phone VARCHAR(20),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS prepared_transactions (
    id SERIAL PRIMARY KEY,
    user_type VARCHAR(50) NOT NULL,
    date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    trans_type VARCHAR(50) NOT NULL,
    amount DECIMAL(15,5) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    status_id INTEGER REFERENCES transaction_statuses(id),
    sender_bank VARCHAR(255),
    receiver_inn VARCHAR(12),
    receiver_phone VARCHAR(20),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS prepared_transactions;
DROP TABLE IF EXISTS transaction_statuses;
DROP TABLE IF EXISTS categories;
DROP TYPE IF EXISTS transaction_type;
-- +goose StatementEnd 