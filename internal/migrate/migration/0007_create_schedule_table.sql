CREATE TABLE schedule (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name text NOT NULL,
    expr text NOT NULL,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    UNIQUE (name, project_id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
