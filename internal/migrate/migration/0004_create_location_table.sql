CREATE TYPE kind AS ENUM ('s3', 'ftp', 'ftps', 'sftp');

CREATE TABLE location (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    kind kind NOT NULL,
    info jsonb NOT NULL,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
