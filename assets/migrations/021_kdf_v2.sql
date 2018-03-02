-- +migrate Up

-- first part is cleanup of old kdf table nobody used
-- first foreign key
ALTER TABLE wallets
  DROP CONSTRAINT wallets_kdf_fkey;
-- then table
DROP TABLE kdf;

-- second part is actual migration
-- create base kdf table
CREATE TABLE kdf (
  version   BIGINT PRIMARY KEY,
  algorithm TEXT   NOT NULL,
  bits      BIGINT NOT NULL,
  n         BIGINT NOT NULL,
  r         BIGINT NOT NULL,
  p         BIGINT NOT NULL
);

-- create many-to-many link
CREATE TABLE kdf_wallets (
  version  BIGINT NOT NULL,
  wallet   TEXT   NOT NULL,
  salt     TEXT   NOT NULL
);

-- insert default kdf versions
INSERT INTO kdf (version, algorithm, bits, n, r, p) VALUES ('1', 'scrypt', 256, 4096, 8, 1);
-- this one is same but version bump signifies update to signup flow for clients
INSERT INTO kdf (version, algorithm, bits, n, r, p) VALUES ('2', 'scrypt', 256, 4096, 8, 1);

-- relate kdf/kdf_wallet/wallets tables
ALTER TABLE kdf_wallets
  ADD CONSTRAINT kdf_wallets_wallets_fkey FOREIGN KEY (wallet) REFERENCES wallets (email) on delete cascade;

ALTER TABLE kdf_wallets
  ADD CONSTRAINT kdf_wallets_kdf_fkey FOREIGN KEY (version) REFERENCES kdf (version);

-- migrate salt to the new table
INSERT INTO kdf_wallets (version, wallet, salt) SELECT
                                                  w.kdf_id,
                                                  w.email,
                                                  w.salt
                                                FROM wallets w;

-- drop now-in-kdf-wallets columns
ALTER TABLE wallets
  DROP COLUMN kdf_id,
  DROP COLUMN salt;

-- +migrate Down

-- we won't even try to revert it to state prior to migration, so just making down pass
ALTER TABLE wallets
  ADD COLUMN kdf_id INT,
  ADD COLUMN salt TEXT;

ALTER TABLE wallets
  ADD CONSTRAINT wallets_kdf_fkey FOREIGN KEY (kdf_id) REFERENCES kdf (version);


update wallets w set kdf_id=(select version from kdf_wallets k where w.email=k.wallet);
update wallets w set salt=(select salt from kdf_wallets k where w.email=k.wallet);
DROP TABLE kdf_wallets;