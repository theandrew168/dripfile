CREATE TABLE history (
    id bigserial PRIMARY KEY,
    transfer_id bigint NOT NULL,  -- not an FK in case transfer gets deleted
    bytes bigint NOT NULL,
    status text NOT NULL,
    started_at timestamptz NOT NULL,
    finished_at timestamptz NOT NULL,
    project_id bigint NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
