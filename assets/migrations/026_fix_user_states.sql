-- +migrate Up

alter table users alter column type drop not null;
alter table users alter column state drop not null;

-- +migrate Down

alter table users alter column type set not null;
alter table users alter column state set not null;
