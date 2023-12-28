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

    CONSTRAINT pk_user_belongs_to_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    CONSTRAINT un_email_is_unique UNIQUE (email)
);

CREATE TABLE monitors (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    account_id VARCHAR(26) NOT NULL,
    endpoint_url VARCHAR(2048) NOT NULL,
    is_endpoint_up BOOLEAN DEFAULT false NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    last_checked_at TIMESTAMPTZ,
    check_interval_in_seconds INT NOT NULL DEFAULT 30,

    CONSTRAINT pk_monitor_belongs_to_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE incidents (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    monitor_id VARCHAR(26) NOT NULL,
    resolved_at TIMESTAMPTZ,
    cause VARCHAR(300),
    response_body TEXT,
    response_headers TEXT,
    response_status SMALLINT NOT NULL,
    request_headers TEXT,
    checked_url VARCHAR(2048) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,

    CONSTRAINT pk_incident_belongs_to_monitor FOREIGN KEY (monitor_id) REFERENCES monitors (id) ON DELETE CASCADE
);

CREATE TYPE INCIDENT_ACTION_TYPE AS ENUM('resolved', 'acknowledged', 'created');

CREATE TABLE incident_actions (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    description VARCHAR(512),
    action_type INCIDENT_ACTION_TYPE DEFAULT 'created' NOT NULL,
    incident_id VARCHAR(26) NOT NULL,
    taken_by_user_with_id VARCHAR(26),
    at TIMESTAMPTZ DEFAULT now() NOT NULL,

    CONSTRAINT pk_incident_action_belongs_to_incidents FOREIGN KEY (incident_id) REFERENCES incidents (id) ON DELETE CASCADE,
    CONSTRAINT pk_incident_action_can_be_taken_by_user FOREIGN KEY (taken_by_user_with_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE subscribers (
    user_id VARCHAR(26) NOT NULL,
    monitor_id VARCHAR(26) NOT NULL,

    CONSTRAINT pk_user_id_monitor_id PRIMARY KEY (user_id, monitor_id),
    CONSTRAINT fk_is_attached_to_a_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_is_attached_to_a_monitor FOREIGN KEY (monitor_id) REFERENCES monitors (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subscribers;
DROP TABLE incident_actions;
DROP TABLE incidents;
DROP TABLE monitors;
DROP TABLE users;
DROP TABLE accounts;
DROP TYPE ROLE_TYPE;
DROP TYPE INCIDENT_ACTION_TYPE;
-- +goose StatementEnd
