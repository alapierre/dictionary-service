alter table dictionary
    add constraint dictionary_dictionary_metadata_type_tenant_fk
        foreign key (type, tenant) references dictionary_metadata;

alter table dictionary
    add lp varchar(64) default '' not null;