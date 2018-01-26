-- +migrate Up

alter table users add column kyc_sequence int not null default 0;
alter table users add column reject_reason text not null default '';

-- +migrate Down

alter table users drop column kyc_sequence;
alter table users drop column reject_reason;
