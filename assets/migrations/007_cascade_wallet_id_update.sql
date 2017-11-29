-- +migrate Up

alter table email_tokens drop constraint email_tokens_wallets_fkey;
alter table email_tokens add constraint email_tokens_wallets_fkey foreign key (wallet_id) references wallets(wallet_id) on delete cascade on update cascade;

alter table tfa_backends drop constraint tfa_backends_wallets_fket;
alter table tfa_backends add constraint tfa_backends_wallets_fkey foreign key (wallet_id) references wallets(wallet_id) on delete cascade on update cascade;

-- +migrate Down

alter table email_tokens drop constraint email_tokens_wallets_fkey;
alter table email_tokens add constraint email_tokens_wallets_fkey foreign key (wallet_id) references wallets(wallet_id) on delete cascade;

alter table tfa_backends drop constraint tfa_backends_wallets_fkey;
alter table tfa_backends add constraint tfa_backends_wallets_fket foreign key (wallet_id) references wallets(wallet_id) on delete cascade;
