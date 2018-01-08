-- +migrate Up

alter table blobs add column relationships jsonb not null default '{}';

-- +migrate Down

alter table blobs drop column relationships;
