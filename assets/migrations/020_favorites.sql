-- +migrate Up

CREATE TABLE favorites (
  id         BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  owner      TEXT      NOT NULL
    CONSTRAINT favorites_users_fkey
    REFERENCES users (address)
    ON UPDATE CASCADE ON DELETE CASCADE,
  type       INT       NOT NULL,
  key        TEXT      NOT NULL
);

CREATE UNIQUE INDEX favorites_unique_per_user
  ON favorites (owner, type, key);

-- +migrate Down

DROP INDEX favorites_unique_per_user;
DROP TABLE favorites;