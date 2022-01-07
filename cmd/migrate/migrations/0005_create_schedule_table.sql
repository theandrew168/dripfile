CREATE TABLE schedule (
    id serial PRIMARY KEY,
    expr text NOT NULL,
    project_id integer NOT NULL REFERENCES project(id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
