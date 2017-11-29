-- +migrate Up

alter table wallets drop column tfa_public_key;
alter table wallets drop column tfa_keychain_data;
alter table wallets drop column tfa_salt;

-- +migrate Down

alter table wallets add column tfa_public_key varchar(64) not null default '';
alter table wallets add column tfa_keychain_data text not null default '';
alter table wallets add column tfa_salt varchar(64) not null default '';
