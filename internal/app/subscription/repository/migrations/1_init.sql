-- +goose Up
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS subscriptions(email TEXT PRIMARY KEY, bundle_id INTEGER, is_deleted BOOL);
CREATE TABLE IF NOT EXISTS balances(email TEXT PRIMARY KEY, current_balance INTEGER);
CREATE TABLE IF NOT EXISTS balances_history(email TEXT, amount INTEGER, operation TEXT, made_at TIMESTAMP WITH TIME ZONE);
COMMIT;

BEGIN TRANSACTION;
CREATE INDEX IF NOT EXISTS idx_subscriptions ON subscriptions USING BTREE (email, bundle_id, is_deleted);
CREATE INDEX IF NOT EXISTS idx_balances ON balances USING BTREE (email, current_balance);
CREATE INDEX IF NOT EXISTS idx_balances_history ON balances_history USING BTREE (email, amount, operation);
COMMIT;