-- +migrate Up

CREATE TABLE airdrops (
  owner TEXT NOT NULL
    CONSTRAINT airdrops_users_fkey
    REFERENCES users (address)
    ON UPDATE CASCADE ON DELETE CASCADE,
  state INT  NOT NULL
);

CREATE UNIQUE INDEX airdrops_unique_per_user
  ON airdrops (owner);

-- +migrate Down

DROP index airdrops_unique_per_user;
DROP TABLE airdrops;


