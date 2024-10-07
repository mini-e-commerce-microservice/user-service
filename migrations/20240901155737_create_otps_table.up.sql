CREATE TABLE otps
(
    id      bigserial primary key,
    user_id bigint       not null
        references users on delete cascade,
    usecase varchar(100) not null,
    code    varchar(50)  not null,
    type    varchar(50)  not null,
    counter smallint     not null,
    expired timestamp    not null,
    token   text
);