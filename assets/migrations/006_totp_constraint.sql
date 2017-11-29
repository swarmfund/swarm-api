-- +migrate Up

create unique index tfa_backends_totp_constraint on tfa_backends(wallet_id, backend) where (backend='totp');

-- +migrate Down

drop index tfa_backends_totp_constraint;