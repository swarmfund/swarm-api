-- +migrate Up

alter table favorites add column email text null;
alter table favorites alter column owner drop not null;
alter table favorites add constraint favorites_owner_or_email check ((email is null) != (owner is null));

drop index favorites_unique_per_user;
create unique index favorites_unique_per_owner
  on favorites (owner, type, key) where owner is not null;
create unique index favorites_unique_per_email
  on favorites (email, type, key) where email is not null;

-- +migrate Down

alter table favorites drop column email;
delete from favorites where owner is null;
alter table favorites alter column owner set not null;

drop index favorites_unique_per_owner;
drop index favorites_unique_per_email;
create unique index favorites_unique_per_user
  on favorites (owner, type, key);
