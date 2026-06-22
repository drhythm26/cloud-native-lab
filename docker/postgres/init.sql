CREATE TABLE IF NOT EXISTS releases (
    id TEXT PRIMARY KEY,
    service_name TEXT NOT NULL,
    version TEXT NOT NULL,
    environment TEXT NOT NULL,
    status TEXT NOT NULL,
    owner TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
