-- +migrate Up

CREATE TABLE authorized_device (
    id bigserial PRIMARY KEY,
    wallet_id integer,
    fingerprint varchar(64) NOT NULL,
    details jsonb DEFAULT '{}'::jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()),
    last_login_at timestamp without time zone DEFAULT timezone('utc'::text, now())
);

CREATE TABLE hmac_keys (
    id bigserial PRIMARY KEY ,
    public varchar(64),
    secret varchar(64)
);

CREATE TABLE kyc_entities (
    id bigserial primary key,
    user_id bigint NOT NULL,
    data jsonb NOT NULL,
    type integer NOT NULL,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now())
);

CREATE TABLE kyc_tracker (
    id bigserial primary key,
    user_id bigint NOT NULL,
    last_state varchar(64) NOT NULL
);

CREATE TABLE notifications (
    id bigserial primary key,
    address varchar(64) NOT NULL,
    email varchar(255),
    type integer NOT NULL
);

CREATE TABLE organization_wallets (
    id bigserial primary key,
    wallet_id bigint NOT NULL,
    organization_address varchar(64),
    operation varchar(20) DEFAULT '0'::varchar NOT NULL
);

CREATE TABLE pending_transaction_signers (
    pending_transaction_id bigint NOT NULL,
    signer_identity bigint NOT NULL,
    signer_public_key varchar(64) NOT NULL,
    signer_name text DEFAULT ''::text NOT NULL,
    id integer NOT NULL
);

CREATE TABLE pending_transactions (
    id bigserial primary key,
    tx_hash varchar(64) NOT NULL,
    tx_envelope text NOT NULL,
    operation_type integer NOT NULL,
    state integer NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    source varchar(64) NOT NULL
);

CREATE TABLE recovery_requests (
    id bigserial primary key,
    wallet_id integer,
    account_id varchar(64) NOT NULL,
    email_token varchar(64) NOT NULL,
    code varchar(64) NOT NULL,
    sent_at timestamp without time zone,
    code_shown_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    username varchar(255),
    recovery_wallet_id varchar(255),
    uploaded_at timestamp without time zone
);

CREATE TABLE tfa (
    id bigserial primary key,
    token varchar(64) NOT NULL,
    verified boolean DEFAULT false NOT NULL,
    otp_data jsonb NOT NULL,
    backend bigint
);

CREATE TABLE tfa_backends (
    id bigserial primary key,
    wallet_id bigint,
    backend integer NOT NULL,
    priority integer,
    details jsonb DEFAULT '{}'::jsonb NOT NULL
);

CREATE TABLE users (
    id bigserial primary key,
    address varchar(64) NOT NULL,
    email varchar(64) NOT NULL,
    user_type varchar(64),
    state varchar(64) NOT NULL,
    recovery_state integer DEFAULT 0 NOT NULL,
    limit_review_state integer DEFAULT 0 NOT NULL,
    documents jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    documents_version bigint DEFAULT 0 NOT NULL
);

CREATE TABLE wallets (
    id bigserial primary key,
    wallet_id varchar(256) NOT NULL,
    username varchar(256) NOT NULL,
    salt varchar(256) NOT NULL,
    kdf_params text NOT NULL,
    keychain_data text,
    verification_token varchar(255) DEFAULT ''::varchar NOT NULL,
    verified boolean DEFAULT false NOT NULL,
    account_id varchar(64) NOT NULL,
    current_account_id varchar(64) NOT NULL,
    tfa_public_key varchar(64) NOT NULL,
    tfa_keychain_data text NOT NULL,
    tfa_salt varchar(64) NOT NULL
);

ALTER TABLE ONLY authorized_device
    ADD CONSTRAINT authorized_device_fingerprint_key UNIQUE (fingerprint);

ALTER TABLE ONLY kyc_tracker
    ADD CONSTRAINT kyc_tracker_user_id_unique UNIQUE (user_id);

ALTER TABLE ONLY notifications
    ADD CONSTRAINT notifications_address_type_unique UNIQUE (address, type);

ALTER TABLE ONLY recovery_requests
    ADD CONSTRAINT recovery_requests_wallet_id_key UNIQUE (wallet_id);

ALTER TABLE ONLY tfa
    ADD CONSTRAINT tfa_token_unique UNIQUE (token);

ALTER TABLE ONLY users
    ADD CONSTRAINT unique_address UNIQUE (address);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_unique_email UNIQUE (email);

ALTER TABLE ONLY wallets
    ADD CONSTRAINT wallets_wallet_id_key UNIQUE (wallet_id);

CREATE INDEX users_by_email ON users USING btree (email);

ALTER TABLE ONLY authorized_device
    ADD CONSTRAINT authorized_device_wallet_id_fkey FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE;

ALTER TABLE ONLY kyc_entities
    ADD CONSTRAINT kyc_entities_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY kyc_tracker
    ADD CONSTRAINT kyc_tracker_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY organization_wallets
    ADD CONSTRAINT organization_wallets_organization_address_fkey FOREIGN KEY (organization_address) REFERENCES users(address) ON DELETE CASCADE;

ALTER TABLE ONLY organization_wallets
    ADD CONSTRAINT organization_wallets_wallet_id_fkey FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE;

ALTER TABLE ONLY pending_transaction_signers
    ADD CONSTRAINT pending_transaction_signers_pending_transaction_id_fkey FOREIGN KEY (pending_transaction_id) REFERENCES pending_transactions(id) ON DELETE CASCADE;

ALTER TABLE ONLY recovery_requests
    ADD CONSTRAINT recovery_requests_wallet_id_fkey FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE;

ALTER TABLE ONLY tfa
    ADD CONSTRAINT tfa_backend_fkey FOREIGN KEY (backend) REFERENCES tfa_backends(id) ON DELETE CASCADE;

ALTER TABLE ONLY tfa_backends
    ADD CONSTRAINT tfa_backends_wallet_id_fkey FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE;

-- +migrate Down

drop table authorized_device           ;
drop table hmac_keys                   ;
drop table kyc_entities                ;
drop table kyc_tracker                 ;
drop table notifications               ;
drop table organization_wallets        ;
drop table pending_transaction_signers ;
drop table pending_transactions        ;
drop table recovery_requests           ;
drop table tfa                         ;
drop table tfa_backends                ;
drop table users                       ;
drop table wallets                     ;