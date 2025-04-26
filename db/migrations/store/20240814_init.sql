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
    created_at    timestamp      not null,
    updated_at    timestamp,
    deleted_at    timestamp
);

CREATE TABLE IF NOT EXISTS cart
(
    id          varchar(20)     primary key,
    product     varchar         not null        UNIQUE,
    quantity    integer         not null
);

CREATE TABLE IF NOT EXISTS users_cart
(
    user_id     varchar(20)      references users (id) on delete cascade,
    cart_id     varchar(20)      references cart (id) on delete cascade,
    primary key (user_id, cart_id)
);

CREATE TABLE IF NOT EXISTS products
(
    id             varchar(20)     primary key,
    product        varchar         not null        UNIQUE,
    description    varchar         not null
);

CREATE TABLE IF NOT EXISTS cart_products
(
    cart_id       varchar(20)    references cart (id) on delete cascade,
    product_id    varchar(20)    references products (id) on delete cascade,
    primary key (cart_id, product_id)
);

CREATE TABLE IF NOT EXISTS referral
(
    id           varchar(20)     primary key,
    short_url    varchar         NOT NULL        UNIQUE,
    counter      integer         DEFAULT(0)
);

CREATE TABLE IF NOT EXISTS users_referral
(
    user_id        varchar(20)    references users (id) on delete cascade,
    referral_id    varchar(20)    references referral (id) on delete cascade,
    primary key (user_id, referral_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd