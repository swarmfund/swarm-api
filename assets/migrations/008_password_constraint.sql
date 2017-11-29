-- +migrate Up

create unique index tfa_backends_password_constraint on tfa_backends(wallet_id, backend) where (backend='password');

-- +migrate Down

drop index tfa_backends_password_constraint;