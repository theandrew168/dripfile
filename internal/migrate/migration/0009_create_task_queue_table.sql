-- https://webapp.io/blog/postgres-is-the-answer/
CREATE TABLE task_queue (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    kind text NOT NULL,
    info jsonb NOT NULL,
    status text NOT NULL,
    error text NOT NULL,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
