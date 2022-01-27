CREATE TABLE transfer (
    id bigserial PRIMARY KEY,
    pattern text NOT NULL,
    src_id bigint NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    dst_id bigint NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    project_id bigint NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
