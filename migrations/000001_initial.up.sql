CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    plaid_item_id TEXT NOT NULL UNIQUE,
    access_token_enc TEXT NOT NULL,
    sync_status TEXT DEFAULT 'active',
    next_cursor TEXT DEFAULT '',
    error_message TEXT,
    last_synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT chk_sync_status CHECK (sync_status IN ('active', 'error', 'resyncing'))
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    plaid_transaction_id TEXT NOT NULL UNIQUE,
    plaid_pending_id TEXT,
    amount_cents INT NOT NULL,
    currency_code TEXT NOT NULL,
    date DATE NOT NULL,
    merchant_name TEXT,
    status TEXT NOT NULL,
    raw_payload JSONB,
    is_removed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT chk_status CHECK (status IN ('pending', 'posted'))
);

CREATE INDEX IF NOT EXISTS idx_transactions_item_id ON transactions(item_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_plaid_id ON transactions(plaid_transaction_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);

CREATE INDEX IF NOT EXISTS idx_items_plaid_id ON items(plaid_item_id);
CREATE INDEX IF NOT EXISTS idx_items_tenant_id ON items(tenant_id);
