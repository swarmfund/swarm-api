-- +migrate Up

truncate table users cascade;
alter table users drop column user_type;
alter table users drop column state;
alter table users add column type int not null;
alter table users add column state int not null;

-- +migrate Down

truncate table users cascade;
alter table users drop column type;
alter table users drop column state;
alter table users add column user_type varchar(64) not null;
alter table users add column state varchar(64) not null;
