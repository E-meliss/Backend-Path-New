-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
  id              BIGSERIAL PRIMARY KEY,
  username        TEXT NOT NULL UNIQUE,
  email           TEXT NOT NULL UNIQUE,
  password_hash   TEXT NOT NULL,
  role            TEXT NOT NULL DEFAULT 'user',
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE balances (
  user_id         BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  amount          NUMERIC(18,2) NOT NULL DEFAULT 0,
  last_updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE transactions (
  id            BIGSERIAL PRIMARY KEY,
  from_user_id  BIGINT REFERENCES users(id) ON DELETE SET NULL,
  to_user_id    BIGINT REFERENCES users(id) ON DELETE SET NULL,
  amount        NUMERIC(18,2) NOT NULL CHECK (amount > 0),
  type          TEXT NOT NULL,
  status        TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_from_user_created_at
  ON transactions(from_user_id, created_at DESC);

CREATE INDEX idx_transactions_to_user_created_at
  ON transactions(to_user_id, created_at DESC);

CREATE INDEX idx_transactions_status_created_at
  ON transactions(status, created_at DESC);

CREATE TABLE audit_logs (
  id          BIGSERIAL PRIMARY KEY,
  entity_type TEXT NOT NULL,
  entity_id   TEXT NOT NULL,
  action      TEXT NOT NULL,
  details     JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_entity_created_at
  ON audit_logs(entity_type, entity_id, created_at DESC);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
