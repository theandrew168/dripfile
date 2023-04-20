CREATE TABLE transfer (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    pattern text NOT NULL,
    from_location_id uuid NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    to_location_id uuid NOT NULL REFERENCES location(id) ON DELETE CASCADE,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
