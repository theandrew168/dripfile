CREATE TABLE location (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    kind text NOT NULL,
    info bytea NOT NULL,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
