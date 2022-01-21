CREATE TABLE session (
    id text PRIMARY KEY,
    expiry timestamptz NOT NULL,
    account_id bigint NOT NULL REFERENCES account(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
