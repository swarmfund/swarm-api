-- +migrate Up

create table kdf (
  id int primary key,
  algorithm varchar(255) not null,
  bits integer not null,
  n decimal not null,
  r integer not null,
  p integer not null
);
insert into kdf (id, algorithm, bits, n, r, p) values (1, 'scrypt', 256, 4096, 8, 1);

alter table wallets rename column username to email;
alter table wallets drop column kdf_params;
alter table wallets add column kdf_id int not null;
alter table wallets add constraint wallets_kdf_fkey foreign key (kdf_id) references kdf(id);
alter table wallets drop column verification_token;

-- +migrate Down

alter table wallets rename column email to username;
alter table wallets add column kdf_params text;
alter table wallets drop column kdf_id;
alter table wallets add column verification_token text;
drop table kdf;
