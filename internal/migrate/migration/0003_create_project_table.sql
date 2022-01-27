CREATE TABLE project (
    id bigserial PRIMARY KEY,
    name text NOT NULL,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
