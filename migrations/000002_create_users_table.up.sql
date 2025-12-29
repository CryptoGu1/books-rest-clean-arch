CREATE TABLE users (
    id serial primary key ,
    name varchar(255) not null ,
    email varchar(255) not null unique ,
    password_hash varchar(255) not null,
    registered_at timestamp default now()
);