CREATE TABLE job_schedule (
    job_id integer NOT NULL REFERENCES job(id),
    schedule_id integer NOT NULL REFERENCES schedule(id),
    UNIQUE (job_id, schedule_id),

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX job_schedule_job_id_idx ON job_schedule(job_id);
CREATE INDEX job_schedule_schedule_id_idx ON job_schedule(schedule_id);
