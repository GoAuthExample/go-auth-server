create database if not exists blueprint;

use blueprint;

create table if not exists users (
    user_id serial PRIMARY KEY,
    auth_id varchar(255) NOT NULL UNIQUE,
    email varchar(255) NOT NULL UNIQUE,
    name varchar(255),
    picture varchar(255)
);