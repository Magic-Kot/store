-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id            varchar(20)    primary key,
    person_id     serial         not null        unique,
    login         varchar(30)    not null        unique,
    password      varchar        not null,
    name          varchar        default(''),
    surname       varchar        default(''),
    email         varchar        default('')     unique,
    age           integer        default(0),
    avatar        varchar        default(''),
    created_at    timestamp      not null
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
    id             SERIAL     PRIMARY KEY,
    product        VARCHAR    NOT NULL        UNIQUE,
    description    VARCHAR    NOT NULL
);

CREATE TABLE IF NOT EXISTS cart_products
(
    id            SERIAL     PRIMARY KEY,
    cart_id       INTEGER references cart (id) on delete cascade    NOT NULL,
    product_id    INTEGER references products (id) on delete cascade    NOT NULL
);

CREATE TABLE IF NOT EXISTS referral
(
    id           SERIAL     PRIMARY KEY,
    short_url    VARCHAR    NOT NULL        UNIQUE,
    counter      INTEGER    DEFAULT(0)
);

CREATE TABLE IF NOT EXISTS users_referral
(
    id             SERIAL     PRIMARY KEY,
    user_id        INTEGER references users (id) on delete cascade    NOT NULL,
    referral_id    INTEGER references referral (id) on delete cascade    NOT NULL
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