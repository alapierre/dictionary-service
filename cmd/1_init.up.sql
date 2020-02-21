create table dictionary
(
	key varchar(128) not null,
	type varchar(256) not null,
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
	tenant varchar(128) default ''::character varying not null,
	parent_key varchar(128) not null,
	content jsonb,
	constraint child_pk
		primary key (key, type, tenant),
	constraint child_parent__fk
		foreign key (parent_key, tenant, type) references dictionary (key, tenant, type)
			on delete cascade
);

create view all_dictionaries as
    select key, type, group_id, tenant, content
        from dictionary
    union
    select child.key, child.type, d.group_id, child.tenant, child.content
        from child join dictionary d on child.parent_key = d.key and child.type = d.type and child.tenant = d.tenant;

