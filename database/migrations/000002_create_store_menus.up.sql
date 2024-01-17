CREATE TABLE store_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    secondary_id SERIAL,
    store_id UUID REFERENCES stores(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    image TEXT DEFAULT NULL,
    price INTEGER NOT NULL,
    ordered_count INTEGER DEFAULT 0,
    price_promo INTEGER DEFAULT NULL,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);