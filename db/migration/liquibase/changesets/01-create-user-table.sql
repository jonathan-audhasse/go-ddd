--liquibase formatted sql

--changeset user:create-user-table
--comment: user table
create table "user" (
    id varchar(255) not null primary key,
    username varchar(255) unique not null,
    password varchar(255) not null
);
--comment: insert user {username: alice, pwd: abc123} and {username: bob, pwd: 123abc}
insert into "user" (id, username, password) 
values 
('usr_00', 'alice', encode(sha256('abc123'::bytea)::bytea, 'base64')),
('usr_01', 'bob', encode(sha256('123abc'::bytea)::bytea, 'base64'));
--rollback drop table "user";
--
