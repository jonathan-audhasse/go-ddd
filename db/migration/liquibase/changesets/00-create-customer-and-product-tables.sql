--liquibase formatted sql

--changeset johnny:customer-table labels:my-label context:init-tables
--comment: customer table
create table customer (
    id varchar(255) not null primary key,
    name varchar(255) not null unique,
    email varchar(255)
)
--rollback drop table customer;

--changeset johnny:product-table labels:my-label context:init-tables
--comment: product table
create table product  (
    id varchar(255) not null primary key,
    name varchar(255) not null unique,
    price int not null default 0
)
--rollback drop table product;
