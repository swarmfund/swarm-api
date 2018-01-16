-- +migrate Up

DROP TABLE pending_transaction_signers;
DROP TABLE pending_transactions;

-- +migrate Down

CREATE TABLE pending_transactions
(
  id             BIGSERIAL   NOT NULL
    CONSTRAINT pending_transactions_pkey
    PRIMARY KEY,
  tx_hash        VARCHAR(64) NOT NULL,
  tx_envelope    TEXT        NOT NULL,
  operation_type INTEGER     NOT NULL,
  state          INTEGER     NOT NULL,
  created_at     TIMESTAMP,
  updated_at     TIMESTAMP,
  source         VARCHAR(64) NOT NULL
);


CREATE TABLE pending_transaction_signers
(
  pending_transaction_id BIGINT                  NOT NULL
    CONSTRAINT pending_transaction_signers_pending_transaction_id_fkey
    REFERENCES pending_transactions
    ON DELETE CASCADE,
  signer_identity        BIGINT                  NOT NULL,
  signer_public_key      VARCHAR(64)             NOT NULL,
  signer_name            TEXT DEFAULT '' :: TEXT NOT NULL,
  id                     INTEGER                 NOT NULL
);

