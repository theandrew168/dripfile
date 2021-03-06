CREATE TABLE transfer (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    pattern text NOT NULL,
    src_id uuid NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    dst_id uuid NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    schedule_id uuid NOT NULL REFERENCES schedule(id) ON DELETE CASCADE,
    project_id uuid NOT NULL REFERENCES project(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
