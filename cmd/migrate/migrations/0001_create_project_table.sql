CREATE TABLE project (
    id serial PRIMARY KEY,
    name text NOT NULL,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
