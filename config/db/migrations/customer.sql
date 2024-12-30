--- New Schema
create schema aion_api

-- auto-generated definition

    create table users
    (
        id         SERIAL not null unique,
        name       varchar(255) not null,
        username   varchar(255) not null unique,
        password   varchar(255) not null,
        email      varchar(255) not null unique,
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp
    );

INSERT INTO aion_api.users (name, username, password, email)
VALUES ('Felipe Bravo', 'lechitz', '123456', 'felipe@gmail.com');

INSERT INTO aion_api.users (name, username, password, email)
VALUES ('Silvana Bravo', 'sil', '123456', 'silbravo@gmail.com');