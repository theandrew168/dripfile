CREATE TABLE itinerary (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    from_location_id uuid NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    to_location_id uuid NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    pattern text NOT NULL,

    -- metadata columns
    created_at timestamptz NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
