DROP TABLE IF EXISTS users, locations, visits;
DROP TYPE IF EXISTS gender;

CREATE TYPE gender AS ENUM ('m', 'f');

CREATE TABLE IF NOT EXISTS users (
    id bigint PRIMARY KEY,
    email varchar(100) UNIQUE NOT NULL,
    first_name varchar(50) NOT NULL,
    last_name varchar(50) NOT NULL,
    gender gender,
    birth_date bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS locations (
    id bigint PRIMARY KEY,
    place text NOT NULL,
    country varchar(50) NOT NULL,
    city varchar(50) NOT NULL,
    distance bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS visits (
    id bigint PRIMARY KEY,
    location bigint REFERENCES locations NOT NULL,
    "user" bigint REFERENCES users NOT NULL,
    visited_at bigint NOT NULL,
    mark integer CHECK (mark BETWEEN 0 AND 5)
);
