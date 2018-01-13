-- +migrate Up

alter table users drop column recovery_state;
drop table recovery_requests;

-- +migrate Down

alter table users add column recovery_state int not null default 0;
create table recovery_requests ();
