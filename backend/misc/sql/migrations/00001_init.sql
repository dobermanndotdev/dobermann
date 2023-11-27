-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TYPE ROLE_TYPE AS ENUM('owner', 'admin', 'writer', 'reader');

CREATE TABLE users (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    first_name VARCHAR(64),
    last_name VARCHAR(64),
    email VARCHAR(128) NOT NULL,
    password VARCHAR(128) NOT NULL,
    role ROLE_TYPE NOT NULL DEFAULT 'admin',
    account_id VARCHAR(26) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,

    CONSTRAINT pk_user_belongs_to_account FOREIGN KEY (account_id) REFERENCES accounts(id),
    CONSTRAINT un_email_is_unique_per_account UNIQUE (email, account_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE accounts;
DROP TYPE ROLE_TYPE;
-- +goose StatementEnd
