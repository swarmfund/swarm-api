-- +migrate Up

CREATE TABLE user_states (
  address    TEXT                        NOT NULL,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  state      INT                         NULL,
  type       INT                         NULL,
  kyc_blob   TEXT                        NULL
);

ALTER TABLE user_states
  ADD CONSTRAINT user_states_users_fkey FOREIGN KEY (address) REFERENCES users (address) ON DELETE CASCADE;

ALTER TABLE user_states
  ADD CONSTRAINT user_states_blobs_fkey FOREIGN KEY (kyc_blob) REFERENCES blobs (id);

CREATE UNIQUE INDEX user_states_user_unique
  ON user_states (address);

-- ALTER TABLE users
--   DROP COLUMN state;
-- ALTER TABLE users
--   DROP COLUMN type;

-- +migrate Down

DROP TABLE user_states;

-- ALTER TABLE users
--   ADD COLUMN state INT NOT NULL DEFAULT 1;
-- ALTER TABLE users
--   ADD COLUMN type INT NOT NULL DEFAULT 1;