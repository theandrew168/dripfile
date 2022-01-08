CREATE TYPE kind AS ENUM ('s3', 'ftp', 'ftps', 'sftp');

CREATE TABLE location (
    id serial PRIMARY KEY,
    kind kind NOT NULL,
    info jsonb NOT NULL,
    project_id integer NOT NULL REFERENCES project(id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
