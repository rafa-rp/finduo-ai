package sqlite

// SchemaSQL contains the DDL to set up the finduo database tables in SQLite.
const SchemaSQL = `
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    salary REAL NOT NULL DEFAULT 0.00
);

CREATE TABLE IF NOT EXISTS expenses (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    amount REAL NOT NULL,
    date TEXT NOT NULL,
    category TEXT NOT NULL,
    payer_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_shared INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS monthly_settlements (
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    is_settled INTEGER NOT NULL DEFAULT 0,
    settled_by_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    PRIMARY KEY (year, month)
);
`
