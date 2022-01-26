CREATE TABLE job (
    transfer_id bigint NOT NULL REFERENCES transfer(id) ON DELETE CASCADE,
    schedule_id bigint NOT NULL REFERENCES schedule(id) ON DELETE CASCADE,
    UNIQUE (transfer_id, schedule_id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX job_transfer_id_idx ON job(transfer_id);
CREATE INDEX job_schedule_id_idx ON job(schedule_id);
