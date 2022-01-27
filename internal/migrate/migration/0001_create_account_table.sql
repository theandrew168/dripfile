CREATE TABLE account (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    email citext NOT NULL UNIQUE,
    username text NOT NULL,
    password text NOT NULL,
    verified boolean NOT NULL,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
