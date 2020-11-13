create table configuration
(
    key varchar(256) not null,
    tenant varchar(128) default ''::character varying not null,
    type varchar(256) not null,
    name varchar(256) not null,
    value varchar(512),
    date_from date not null,
    date_to date not null,
    constraint configuration_pk
        primary key (key, tenant, date_from)
);

create index configuration_date_from_date_to_key_tenant_index
    on configuration (date_from, date_to, key, tenant);
