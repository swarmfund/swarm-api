-- +migrate Up

create table email_tokens (
  id bigserial primary key,
  wallet_id varchar(256) not null,
  token varchar(64),
  last_sent_at timestamp without time zone,
  confirmed bool default false
);
alter table email_tokens add constraint email_tokens_wallets_fkey foreign key (wallet_id) references wallets(wallet_id) on delete cascade;
alter table email_tokens add constraint email_tokens_wallet_id_unique unique(wallet_id);
alter table wallets drop column verified;

-- +migrate Down

drop table email_tokens;
alter table wallets add column verified bool not null default false;
