-- https://webapp.io/blog/postgres-is-the-answer/
CREATE TYPE transfer_status AS ENUM ('new', 'running', 'success', 'error');

CREATE TABLE transfer_queue (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    status transfer_status NOT NULL,
    transfer_id uuid NOT NULL REFERENCES transfer(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
