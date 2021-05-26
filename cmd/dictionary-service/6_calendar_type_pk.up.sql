alter table calendar drop constraint calendar_pk;

alter table calendar
    add constraint calendar_pk
        primary key (day, tenant, type);