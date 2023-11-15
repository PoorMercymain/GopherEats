-- +goose Up
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS ingredients
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) UNIQUE,
    unit       VARCHAR(20),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dishes
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) UNIQUE,
    description       VARCHAR(200),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS resources
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) UNIQUE,
    filename   VARCHAR(200),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bundles
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(200) UNIQUE,
    price   INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dishes_ingredients
(
    id            SERIAL PRIMARY KEY,
    dish_id       INT,
    ingredient_id INT,
    qty           FLOAT
);

CREATE TABLE IF NOT EXISTS dishes_resources
(
    id          SERIAL PRIMARY KEY,
    dish_id     INT,
    resource_id INT
);

CREATE TABLE IF NOT EXISTS bundles_dishes
(
    id          SERIAL PRIMARY KEY,
    week_number INT,
    bundle_id   INT,
    dish_id     INT,
    created_at  TIMESTAMP,
    updated_at  TIMESTAMP
);

COMMIT;

-- +goose Down
BEGIN TRANSACTION;

DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS dishes;
DROP TABLE IF EXISTS bundles;
DROP TABLE IF EXISTS resources;

DROP TABLE IF EXISTS dishes_ingredients;
DROP TABLE IF EXISTS dishes_resources;
DROP TABLE IF EXISTS bundles_dishes;

COMMIT;