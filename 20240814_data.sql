-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL         PRIMARY KEY,
    username      VARCHAR(30)    NOT NULL        UNIQUE,
    password      VARCHAR        NOT NULL,
    name          VARCHAR        DEFAULT(''),
    surname       VARCHAR        DEFAULT(''),
    email         VARCHAR        DEFAULT('')     UNIQUE,
    age           INTEGER        DEFAULT(0),
    avatar        VARCHAR        DEFAULT('')
);

CREATE TABLE IF NOT EXISTS sessions
(
    id            SERIAL         PRIMARY KEY,
    userId        INTEGER references users (id) on delete cascade    NOT NULL,
    guid          VARCHAR        DEFAULT(''),
    refreshToken  VARCHAR        DEFAULT(''),
    expiresAt     VARCHAR        DEFAULT('')
);

CREATE TABLE IF NOT EXISTS cart
(
    id          SERIAL     PRIMARY KEY,
    product     VARCHAR    NOT NULL        UNIQUE,
    quantity    INTEGER    NOT NULL
);

CREATE TABLE IF NOT EXISTS users_cart
(
    id          SERIAL     PRIMARY KEY,
    user_id     INTEGER references users (id) on delete cascade    NOT NULL,
    cart_id     INTEGER references cart (id) on delete cascade    NOT NULL
);

CREATE TABLE IF NOT EXISTS products
(
    id          SERIAL     PRIMARY KEY,
    product     VARCHAR    NOT NULL        UNIQUE,
    description VARCHAR    NOT NULL
);

CREATE TABLE IF NOT EXISTS cart_products
(
    id          SERIAL     PRIMARY KEY,
    cart_id     INTEGER references cart (id) on delete cascade    NOT NULL,
    product_id     INTEGER references products (id) on delete cascade    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS users_cart;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS cart_products;
-- +goose StatementEnd