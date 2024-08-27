CREATE TABLE users
(
    id         bigserial primary key,
    email      varchar(255) not null,
    password   varchar(255) not null,
    created_at timestamp    not null,
    updated_at timestamp    not null,
    deleted_at timestamp
);