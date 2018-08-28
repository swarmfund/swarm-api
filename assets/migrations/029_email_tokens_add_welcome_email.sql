-- +migrate Up

alter table email_tokens add column sent_welcome_email bool default false;

-- +migrate Down

alter table email_tokens drop column sent_welcome_email;
