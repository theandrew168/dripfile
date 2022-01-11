CREATE TABLE transfer_schedule (
    transfer_id integer NOT NULL REFERENCES transfer(id) ON DELETE CASCADE,
    schedule_id integer NOT NULL REFERENCES schedule(id) ON DELETE CASCADE,
    UNIQUE (transfer_id, schedule_id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX transfer_schedule_transfer_id_idx ON transfer_schedule(transfer_id);
CREATE INDEX transfer_schedule_schedule_id_idx ON transfer_schedule(schedule_id);
