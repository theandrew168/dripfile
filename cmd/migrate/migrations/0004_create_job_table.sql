CREATE TABLE job (
    id serial PRIMARY KEY,
    pattern text NOT NULL,
    src_id integer NOT NULL REFERENCES location(id),
    dst_id integer NOT NULL REFERENCES location(id),
    project_id integer NOT NULL REFERENCES project(id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
