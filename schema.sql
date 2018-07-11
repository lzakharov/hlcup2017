CREATE TYPE gender AS ENUM ('m', 'f');

CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY,
    email varchar(100) UNIQUE NOT NULL,
    first_name varchar(50) NOT NULL,
    last_name varchar(50) NOT NULL,
    gender gender,
    birth_date integer NOT NULL
);

CREATE TABLE IF NOT EXISTS locations (
    id integer PRIMARY KEY,
    place text NOT NULL,
    country varchar(50) NOT NULL,
    city varchar(50) NOT NULL,
    distance integer NOT NULL
);

CREATE TABLE IF NOT EXISTS visits (
    id integer PRIMARY KEY,
    location integer REFERENCES locations NOT NULL,
    "user" integer REFERENCES users NOT NULL,
    visited_at integer NOT NULL,
    mark integer CHECK (mark BETWEEN 0 AND 5)
);
