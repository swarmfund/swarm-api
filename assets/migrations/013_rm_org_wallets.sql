-- +migrate Up

drop table organization_wallets;


-- +migrate Down

CREATE TABLE organization_wallets
(
  id                   BIGSERIAL                                    NOT NULL
    CONSTRAINT organization_wallets_pkey
    PRIMARY KEY,
  wallet_id            BIGINT                                       NOT NULL
    CONSTRAINT organization_wallets_wallet_id_fkey
    REFERENCES wallets
    ON DELETE CASCADE,
  organization_address VARCHAR(64)
    CONSTRAINT organization_wallets_organization_address_fkey
    REFERENCES users (address)
    ON DELETE CASCADE,
  operation            VARCHAR(20) DEFAULT '0' :: CHARACTER VARYING NOT NULL
);

