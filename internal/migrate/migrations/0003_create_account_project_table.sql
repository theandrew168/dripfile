CREATE TABLE account_project (
    account_id integer NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    project_id integer NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    UNIQUE (account_id, project_id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX account_project_account_id_idx ON account_project(account_id);
CREATE INDEX account_project_project_id_idx ON account_project(project_id);
