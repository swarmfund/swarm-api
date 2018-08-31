-- +migrate Up
drop table tracking;

CREATE TABLE audit_logs (
  id bigserial primary key,
  action int,
  details jsonb,
  performed_at timestamp without time zone not null,
  user_address text NOT NULL
);

ALTER TABLE ONLY audit_logs
  ADD CONSTRAINT audit_users_fkey foreign key (user_address) references users(address) on delete cascade;

-- +migrate Down
drop table audit_logs;

create table tracking (
  id bigserial primary key,
  address text,
  signer text,
  details jsonb
);