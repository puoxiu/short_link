create database if not exists short_link_db;

use short_link_db;
create table if not exists urls (
    id int primary key auto_increment,
    long_link varchar(255) not null,
    sort_code varchar(10) not null unique,
    is_custom tinyint default 0,
    expire_time datetime not null,
    create_time datetime not null default current_timestamp
);

create index idx_sort_code on urls(sort_code);
create index idx_expire_time on urls(expire_time);
