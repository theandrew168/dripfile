CREATE TABLE history (
    id serial PRIMARY KEY,
    job_id integer NOT NULL,  -- not an FK in case job gets deleted
    bytes bigint NOT NULL,
    status text NOT NULL,
    started_at timestamptz NOT NULL,
    finished_at timestamptz NOT NULL,
    project_id integer NOT NULL REFERENCES project(id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
