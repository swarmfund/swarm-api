-- +migrate Up

create table tracking (
  id primary key,
  address text,
  signer text,
  details jsonb
);


-- +migrate Down

drop table tracking;