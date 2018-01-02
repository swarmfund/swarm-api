-- +migrate Up

TRUNCATE kyc_entities;

ALTER TABLE kyc_entities
  DROP COLUMN id;

ALTER TABLE kyc_entities
  ADD COLUMN id CHAR(26) PRIMARY KEY;

CREATE UNIQUE INDEX kyc_entities_individual_constraint
  ON kyc_entities (user_id, type)
  WHERE (type = 1);


-- +migrate Down

TRUNCATE TABLE kyc_entities;

ALTER TABLE kyc_entities
  DROP COLUMN id;

ALTER TABLE kyc_entities
  ADD COLUMN id BIGSERIAL PRIMARY KEY;

DROP INDEX kyc_entities_individual_constraint;
