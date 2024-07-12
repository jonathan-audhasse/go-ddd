--liquibase formatted sql

--changeset johnny:create-healthcheck-table
--comment: healthcheck table
create table healthcheck  (value varchar(255) not null unique)
--rollback drop table healthcheck;
