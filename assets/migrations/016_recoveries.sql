-- +migrate Up

CREATE UNIQUE INDEX wallets_email_unique_constraint ON wallets(email);

CREATE TABLE recoveries
(
  wallet        TEXT NOT NULL
    CONSTRAINT recoveries_wallets_fkey
    REFERENCES wallets (email)
    ON DELETE CASCADE,
  salt          TEXT NOT NULL,
  keychain_data TEXT NOT NULL,
  wallet_id     TEXT NOT NULL,
  address       TEXT NOT NULL
);

CREATE UNIQUE INDEX recoveries_wallet_id_unique_constraint
  ON recoveries (wallet_id);


-- +migrate Down

drop index recoveries_wallet_id_unique_constraint;
drop table recoveries;
drop index wallets_email_unique_constraint;