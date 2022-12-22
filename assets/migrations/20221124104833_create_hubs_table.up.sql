create table hubs
(
    id         int primary key generated always as identity,
    name       varchar(255)                   not null default '',
    location   varchar(255)                   not null default '',
    created_at timestamp(0) without time zone not null default now()::timestamp without time zone,
    updated_at timestamp(0) without time zone not null default '0001-01-01 00:00:00'::timestamp without time zone
);
