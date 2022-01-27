CREATE TYPE role AS ENUM ('owner', 'admin', 'editor', 'viewer');

CREATE TABLE member (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    role role NOT NULL,
    account_id uuid NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    UNIQUE (account_id, project_id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX member_account_id_idx ON member(account_id);
CREATE INDEX member_project_id_idx ON member(project_id);
