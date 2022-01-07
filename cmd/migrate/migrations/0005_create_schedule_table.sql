CREATE TABLE schedule (
    id serial PRIMARY KEY,
    name text NOT NULL,
    expr text NOT NULL,
    project_id integer NOT NULL REFERENCES project(id),
    UNIQUE (name, project_id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
