create table dictionary
(
	key varchar(128) not null,
	type varchar(256) not null,
	group_id varchar(64),
	tenant varchar(128) default ''::character varying not null,
	content jsonb,
	name varchar(256),
	parent_key varchar(128),
	constraint dictionary_pk
		primary key (key, type, tenant),
	constraint dictionary_dictionary_key_type_fk
		foreign key (parent_key, type, tenant) references dictionary
			on delete cascade
);

create or replace view all_dictionaries(key, type, name, group_id, tenant, content, parent, parent_key) as
SELECT dictionary.key,
       dictionary.type,
       dictionary.name,
       dictionary.group_id,
       dictionary.tenant,
       dictionary.content,
       true                    AS parent,
       NULL::character varying AS parent_key
FROM dictionary where dictionary.parent_key is null
UNION
SELECT child.key,
       child.type,
       child.name,
       d.group_id,
       child.tenant,
       child.content,
       false AS parent,
       d.key AS parent_key
FROM dictionary child
         JOIN dictionary d ON child.parent_key = d.key AND child.type = d.type AND
                              child.tenant = d.tenant;


create table translation
(
	key varchar(128) not null,
	type varchar(256) not null,
	tenant varchar(128) default '' not null,
	language char(2) not null,
	name varchar(256) not null,
	constraint translation_pk
		primary key (key, type, tenant, language)
);

alter table translation
	add constraint translation_dictionary_key_type_tenant_fk
		foreign key (key, type, tenant) references dictionary;

create table dictionary_metadata
(
	type varchar(256) not null,
	tenant varchar(128) default ''::character varying not null,
	content jsonb,
	constraint dictionary_metadata_pk
		primary key (type, tenant)
);
