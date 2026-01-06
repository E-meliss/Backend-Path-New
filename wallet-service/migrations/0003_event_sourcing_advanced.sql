-- +goose Up
-- +goose StatementBegin

-- 1) Add multi-currency support and integer cents columns (keep old numeric columns for backward compatibility)
ALTER TABLE balances
  ADD COLUMN IF NOT EXISTS currency TEXT NOT NULL DEFAULT 'USD',
  ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL DEFAULT 0;

-- If amount_cents is zero but numeric amount is non-zero, backfill best-effort
UPDATE balances
SET amount_cents = COALESCE(amount_cents, 0) + (amount * 100)::BIGINT
WHERE amount IS NOT NULL AND amount_cents = 0;

-- balances primary key becomes (user_id, currency) if not already.
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE table_name='balances' AND constraint_type='PRIMARY KEY'
  ) THEN
    -- drop existing PK if it is only user_id
    BEGIN
      ALTER TABLE balances DROP CONSTRAINT balances_pkey;
    EXCEPTION WHEN undefined_object THEN
      -- ignore
    END;
  END IF;
END $$;

ALTER TABLE balances
  ADD PRIMARY KEY (user_id, currency);

ALTER TABLE transactions
  ADD COLUMN IF NOT EXISTS currency TEXT NOT NULL DEFAULT 'USD',
  ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL DEFAULT 0;

UPDATE transactions
SET amount_cents = (amount * 100)::BIGINT
WHERE amount IS NOT NULL AND amount_cents = 0;

CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_from_user ON transactions(from_user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_user ON transactions(to_user_id);

-- 2) Event store
CREATE TABLE IF NOT EXISTS events (
  id          BIGSERIAL PRIMARY KEY,
  stream      TEXT NOT NULL,
  type        TEXT NOT NULL,
  entity_id   TEXT NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  data        JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_stream_id ON events(stream, id);

-- 3) Scheduled transactions
CREATE TABLE IF NOT EXISTS scheduled_transactions (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type         TEXT NOT NULL, -- credit/debit/transfer
  to_user_id   BIGINT REFERENCES users(id) ON DELETE SET NULL,
  currency     TEXT NOT NULL DEFAULT 'USD',
  amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
  run_at       TIMESTAMPTZ NOT NULL,
  status       TEXT NOT NULL DEFAULT 'scheduled', -- scheduled, processed, failed, cancelled
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_scheduled_run_at ON scheduled_transactions(run_at, status);

-- 4) Currency rates (optional for conversion)
CREATE TABLE IF NOT EXISTS currency_rates (
  base_currency   TEXT NOT NULL,
  quote_currency  TEXT NOT NULL,
  rate_ppm        BIGINT NOT NULL, -- rate * 1_000_000
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (base_currency, quote_currency)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS currency_rates;
DROP TABLE IF EXISTS scheduled_transactions;
DROP TABLE IF EXISTS events;

ALTER TABLE transactions
  DROP COLUMN IF EXISTS currency,
  DROP COLUMN IF EXISTS amount_cents;

ALTER TABLE balances
  DROP CONSTRAINT IF EXISTS balances_pkey;

-- try to restore old PK
ALTER TABLE balances
  ADD PRIMARY KEY (user_id);

ALTER TABLE balances
  DROP COLUMN IF EXISTS currency,
  DROP COLUMN IF EXISTS amount_cents;

-- +goose StatementEnd
