CREATE TABLE project (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    billing_id text NOT NULL UNIQUE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
