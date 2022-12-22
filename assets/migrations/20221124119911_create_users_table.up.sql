create table users
(
    id         int primary key generated always as identity,
    name       varchar(255)                   not null default '',
    title      varchar(255)                   not null default '',
    team_id    int                            not null default 0,
    created_at timestamp(0) without time zone not null default now()::timestamp without time zone,
    updated_at timestamp(0) without time zone not null default '0001-01-01 00:00:00'::timestamp without time zone
);
