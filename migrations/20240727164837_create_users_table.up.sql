CREATE TABLE users
(
    id                bigserial primary key,
    email             varchar(255) not null,
    password          varchar(255) not null,
    is_email_verified boolean      not null default false,
    created_at        TIMESTAMPTZ    not null,
    updated_at        TIMESTAMPTZ    not null,
    deleted_at        TIMESTAMPTZ,
    trace_parent       varchar
);