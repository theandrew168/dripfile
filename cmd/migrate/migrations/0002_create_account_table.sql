CREATE TYPE role AS ENUM ('owner', 'admin', 'editor', 'viewer');

CREATE TABLE account (
    id serial PRIMARY KEY,
    email citext NOT NULL UNIQUE,
    password text NOT NULL,
    verified boolean NOT NULL,
    role role NOT NULL,
    project_id integer NOT NULL REFERENCES project(id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
