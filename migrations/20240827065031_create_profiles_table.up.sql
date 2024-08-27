CREATE TABLE profiles
(
    id               bigserial primary key,
    user_id          bigint       not null
        references users on delete cascade,
    full_name        varchar(255) not null,
    image_profile    varchar(255),
    background_image varchar(255),
    created_at       timestamp    not null,
    updated_at       timestamp    not null,
    deleted_at       timestamp
);