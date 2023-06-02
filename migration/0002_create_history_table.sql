CREATE TABLE history (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    total_bytes bigint NOT NULL,
    started_at timestamptz NOT NULL,
    finished_at timestamptz NOT NULL,
    transfer_id uuid REFERENCES transfer(id) ON DELETE SET NULL,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
