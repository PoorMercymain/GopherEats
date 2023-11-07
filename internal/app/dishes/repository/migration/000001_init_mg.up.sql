BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS ingredients
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) UNIQUE,
    unit       VARCHAR(20),
    image_id   INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dishes
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) UNIQUE,
    desc       VARCHAR(1000),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS resources
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) UNIQUE,
    filename   VARCHAR(1000),
    type       VARCHAR(30),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bundles
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(200) UNIQUE,
    desc       VARCHAR(1000),
    image_id   INT,
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

CREATE TABLE IF NOT EXISTS subscription_orders
(
    id          SERIAL PRIMARY KEY,
    client_name VARCHAR(100),
    address     VARCHAR(100),
    week_number INT,
    bundle_id   INT,
    status      VARCHAR(100),
    created_at  TIMESTAMP,
    updated_at  TIMESTAMP
);


COMMIT;