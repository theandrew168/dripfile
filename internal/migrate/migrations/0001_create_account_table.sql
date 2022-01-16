CREATE TYPE role AS ENUM ('owner', 'admin', 'editor', 'viewer');

CREATE TABLE account (
    id bigserial PRIMARY KEY,
    email citext NOT NULL UNIQUE,
    username text NOT NULL,
    password text NOT NULL,
    verified boolean NOT NULL,
    role role NOT NULL,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
