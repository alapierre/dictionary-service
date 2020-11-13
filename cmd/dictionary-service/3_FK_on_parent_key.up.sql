alter table dictionary
    add constraint dictionary_dictionary_key_tenant_fk
        foreign key (parent_key, tenant, type) references dictionary.dictionary (key, tenant, type);
