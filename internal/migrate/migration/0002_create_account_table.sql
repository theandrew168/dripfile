CREATE TYPE role AS ENUM ('owner', 'admin', 'editor', 'viewer');

CREATE TABLE account (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email citext NOT NULL UNIQUE,
    username text NOT NULL,
    password text NOT NULL,
    role role NOT NULL,
    verified boolean NOT NULL,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
