--liquibase formatted sql

--changeset johnny:create-healthcheck-table
--comment: healthcheck table
create table healthcheck  (value varchar(255) not null unique);
insert into healthcheck (value) values ('J0n@th@n');
--rollback drop table healthcheck;
