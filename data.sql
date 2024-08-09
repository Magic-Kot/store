CREATE TABLE IF NOT EXISTS users
(
    id          SERIAL  PRIMARY KEY,
    login       TEXT    NOT NULL        UNIQUE,
    password    TEXT    NOT NULL
);
--CREATE INDEX IF NOT EXISTS idx_login ON users(login);

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