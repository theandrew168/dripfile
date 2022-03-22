CREATE TABLE history (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    bytes bigint NOT NULL,
    status text NOT NULL,
    started_at timestamptz NOT NULL,
    finished_at timestamptz NOT NULL,
    transfer_id uuid NOT NULL,  -- not an FK in case transfer gets deleted
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
