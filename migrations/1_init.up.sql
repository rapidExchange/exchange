CREATE TABLE IF NOT EXISTS users (
    id int AUTO_INCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash VARCHAR(512) NOT NULL,
    PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS balance (
    email TEXT NOT NULL,
    ticker TEXT NOT NULL,
    quantity DECIMAL(17, 8)
);