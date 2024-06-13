--liquibase formatted sql

--changeset johnny:order-table labels:another-one
--comment: order table
create table "order" (
    id varchar(255) not null primary key,
    customer_id varchar(255) not null references customer(id) on delete cascade,
    created date not null default now(),
    last_modified date not null default now()
)
--rollback drop table "order";

--changeset johnny:order_item-table
--comment: order_item table
create table order_item (
    id varchar(255) not null primary key,
    order_id varchar(255) references "order"(id) on delete cascade,
    product_id varchar(255) references product(id) on delete cascade,
    qty int not null default 1,
    discount float not null check (discount >= 0 and discount <= 1) default 0
)
--rollback drop table order_item;
