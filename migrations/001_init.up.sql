CREATE TABLE listings (
    uuid TEXT PRIMARY KEY,
    source TEXT NOT NULL DEFAULT 'realt.by',
    code BIGINT,
    category INT,

    price NUMERIC(14,2),
    price_currency INT,
    price_per_m2 NUMERIC(14,2),

    rooms INT,
    area_total NUMERIC(8,2),
    area_living NUMERIC(8,2),
    area_kitchen NUMERIC(8,2),
    storey INT,
    storeys INT,
    building_year INT,

    town_name TEXT,
    district_name TEXT,
    street_name TEXT,
    house_number INT,
    address TEXT,

    title TEXT,
    description TEXT,
    images TEXT[],

    source_created_at TIMESTAMPTZ,
    source_updated_at TIMESTAMPTZ,
    raise_date TIMESTAMPTZ,

    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    notified BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_listings_district ON listings (district_name);
CREATE INDEX idx_listings_price ON listings (price_currency, price);
CREATE INDEX idx_listings_created ON listings (source_created_at DESC);

CREATE TABLE price_history (
    id BIGSERIAL PRIMARY KEY,
    uuid TEXT NOT NULL REFERENCES listings(uuid) ON DELETE CASCADE,
    price NUMERIC(14,2),
    currency INT,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_price_history_uuid ON price_history (uuid);