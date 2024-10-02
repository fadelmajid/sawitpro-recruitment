-- database.sql

CREATE TABLE IF NOT EXISTS estates (
    id UUID PRIMARY KEY,
    width INT NOT NULL,
    length INT NOT NULL
);

CREATE TABLE IF NOT EXISTS trees (
    id UUID PRIMARY KEY,
    estate_id UUID REFERENCES estates(id),
    x INT NOT NULL,
    y INT NOT NULL,
    height INT NOT NULL
);