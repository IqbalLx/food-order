CREATE TABLE stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    secondary_id SERIAL,
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    image TEXT NOT NULL,
    short_desc TEXT,
    "desc" TEXT,
    rating INTEGER,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);