CREATE TABLE session (
    hash text PRIMARY KEY,
    expiry timestamptz NOT NULL,
    account_id uuid NOT NULL REFERENCES account(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX session_expiry_idx ON session(expiry);
