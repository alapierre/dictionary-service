create table dictionary
(
	key varchar(128) not null,
	type varchar(256) not null,
	name varchar(256) not null,
	group_id varchar(64),
	tenant varchar(128) default ''::character varying not null,
	content jsonb,
	constraint dictionary_pk
		primary key (key, type, tenant)
);

create table child
(
	key varchar(128) not null,
	type varchar(256) not null,
	name varchar(256) not null,
	tenant varchar(128) default ''::character varying not null,
	parent_key varchar(128) not null,
	content jsonb,
	constraint child_pk
		primary key (key, type, tenant),
	constraint child_parent__fk
		foreign key (parent_key, tenant, type) references dictionary (key, tenant, type)
			on delete cascade
);

create or replace view all_dictionaries as select key, type, name, group_id, tenant, content, true as "parent", null as "parent_key"
                                           from dictionary
                                           union
                                           select child.key, child.type, child.name, d.group_id, child.tenant, child.content, false, d.key as "parent_key"
                                           from child join dictionary d on child.parent_key = d.key and child.type = d.type and child.tenant = d.tenant;
