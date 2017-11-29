-- +migrate Up
alter table tfa_backends drop column wallet_id;
alter table tfa_backends add column wallet_id varchar(256) not null;
alter table tfa_backends add constraint tfa_backends_wallets_fket foreign key (wallet_id) references wallets(wallet_id) on delete cascade;
alter table tfa_backends alter column priority set not null;
alter table tfa_backends drop column backend;
alter table tfa_backends add column backend varchar(255);

-- +migrate Down
alter table tfa_backends drop column wallet_id;
alter table tfa_backends add column wallet_id bigint;
alter table tfa_backends alter column priority drop not null;
alter table tfa_backends drop column backend;
alter table tfa_backends add column backend int not null;