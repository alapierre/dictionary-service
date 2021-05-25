CREATE EXTENSION if not exists hstore;

create table calendar_type
(
    type varchar(256) not null,
    name varchar(521) not null,
    tenant varchar(128) not null default ''
);

alter table calendar_type
    add constraint calendar_type_pk
        primary key (type, tenant);

create table calendar
(
    tenant varchar(128) not null default '',
    type varchar(256) not null,
    name varchar(512),
    kind varchar(128),
    labels hstore,
    day date not null,
    constraint calendar_type_fk foreign key (type, tenant) references calendar_type(type, tenant)
);

alter table calendar
    add constraint calendar_pk
        primary key (day, tenant);
