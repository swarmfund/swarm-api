-- +migrate Up

CREATE TABLE blobs (
  id            CHAR(52) PRIMARY KEY,
  owner_address CHAR(64) NOT NULL
    CONSTRAINT blobs_users_fkey
    REFERENCES users (address)
    ON UPDATE CASCADE ON DELETE CASCADE,
  type          INT      NOT NULL,
  value         TEXT     NOT NULL
);

-- +migrate Down

DROP TABLE blobs;