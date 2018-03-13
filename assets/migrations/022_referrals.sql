-- +migrate Up

create table referrals (
  -- account id of one who referred
  referrer text not null,
  -- account id of one who have been referred
  referral text not null
);

-- unique constraint for account_id for foreign keys.
-- if you bumped into it implementing multi-sign, sorry, find your way around.
CREATE UNIQUE INDEX wallets_account_id_unique
  ON wallets (account_id);

ALTER TABLE referrals
  ADD CONSTRAINT referrals_referrer_wallets_fkey FOREIGN KEY (referrer) REFERENCES wallets (account_id) on delete cascade;

ALTER TABLE referrals
  ADD CONSTRAINT referrals_referral_wallets_fkey FOREIGN KEY (referral) REFERENCES wallets (account_id) on delete cascade;

CREATE UNIQUE INDEX referrals_referral_unique
  ON referrals (referral);

-- +migrate Down

drop table referrals;
drop index wallets_account_id_unique;
