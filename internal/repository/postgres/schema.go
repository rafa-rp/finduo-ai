package postgres

// SchemaSQL contains the DDL to set up the finduo database tables.
const SchemaSQL = `
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    salary NUMERIC(12, 2) NOT NULL DEFAULT 0.00
);

CREATE TABLE IF NOT EXISTS expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description TEXT NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    date DATE NOT NULL,
    category TEXT NOT NULL,
    payer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_shared BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS monthly_settlements (
    year INT NOT NULL,
    month INT NOT NULL,
    is_settled BOOLEAN NOT NULL DEFAULT FALSE,
    settled_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
    PRIMARY KEY (year, month)
);
`
