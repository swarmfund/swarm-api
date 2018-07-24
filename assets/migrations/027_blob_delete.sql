-- +migrate Up

alter table blobs add column deleted_at timestamp without time zone;

-- +migrate Down

alter table blobs drop column deleted_at;