-- +migrate Up

alter table blobs alter column owner_address drop not null;

-- +migrate Down

alter table blobs alter column owner_address set not null;