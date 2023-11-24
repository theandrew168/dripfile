CREATE TABLE transfer (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    itinerary_id uuid REFERENCES itinerary(id) ON DELETE SET NULL,
    status text NOT NULL,
    progress integer NOT NULL,
    error text NOT NULL,

    -- metadata columns
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
