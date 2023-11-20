-- +goose Up
-- +goose StatementBegin
CREATE TYPE REQUEST_METHOD AS ENUM (
    'GET',
    'HEAD'
    'POST',
    'PUT',
    'DELETE',
    'CONNECT',
    'OPTIONS',
    'TRACE',
    'PATCH'
);

CREATE TYPE CHECK_STATUS AS ENUM (
    'pending',
    'enqueued',
    'checked'
);

CREATE TABLE accounts (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE TABLE teams (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    account_id VARCHAR(26) NOT NULL,
    name VARCHAR(128) NOT NULL,

    CONSTRAINT fk_team_belongs_to_an_account FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE TABLE users (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    account_id VARCHAR(26) NOT NULL,

    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email VARCHAR(250) NOT NULL,

    primary_phone_number VARCHAR(15) NOT NULL,
    secondary_phone_number VARCHAR(15) NOT NULL,
    avatar_url VARCHAR(515),
    timezone VARCHAR(2), --TO BE REVIEWED
    on_holidays_until TIMESTAMPTZ,

    confirmation_code VARCHAR(26),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_member_belongs_to_an_account FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE TABLE team_members (
    team_id VARCHAR(26) NOT NULL,
    member_id VARCHAR(26) NOT NULL,

    CONSTRAINT pk_team_id_member_id PRIMARY KEY (team_id, member_id),
    CONSTRAINT fk_team_has_members FOREIGN KEY (team_id) REFERENCES teams(id),
    CONSTRAINT fk_member_belongs_to_team FOREIGN KEY (member_id) REFERENCES users(id)
);

CREATE TABLE monitors (
  id VARCHAR(26) NOT NULL PRIMARY KEY,
  endpoint VARCHAR(255) NOT NULL,

  -- settings
  recovered_only_after        SMALLINT NOT NULL,
  start_an_incident_after        SMALLINT NOT NULL,
  check_interval               SMALLINT NOT NULL,
  alert_domain_expiration_within SMALLINT NOT NULL,

  -- ssl verification
  ssl_verification_enabled BOOLEAN DEFAULT FALSE NOT NULL,
  verify_ssl_expiration_within SMALLINT,

  -- request params
  request_method REQUEST_METHOD NOT NULL DEFAULT 'GET',
  request_timeout INTEGER NOT NULL,
  request_body TEXT,
  follow_redirects BOOLEAN NOT NULL DEFAULT FALSE,
  keep_cookies_while_redirecting BOOLEAN NOT NULL DEFAULT FALSE,
  expected_response_status INTEGER NOT NULL,

  -- http auth
  basic_auth_username VARCHAR(64),
  basic_auth_password VARCHAR(128),

  -- maintenance
  maintenance_from TIMESTAMPTZ,
  maintenance_to TIMESTAMPTZ,

  is_up BOOLEAN DEFAULT false NOT NULL,
  is_paused BOOLEAN DEFAULT false NOT NULL,

  account_id VARCHAR(26) NOT NULL,
  team_id VARCHAR(26) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
  last_checked_at TIMESTAMPTZ DEFAULT now() NOT NULL,
  check_status CHECK_STATUS DEFAULT 'enqueued' NOT NULL,

  CONSTRAINT fk_monitor_belongs_to_account FOREIGN KEY (account_id) REFERENCES accounts (id),
  CONSTRAINT fk_monitor_belongs_to_team FOREIGN KEY (team_id) REFERENCES teams (id)
);

CREATE TABLE monitor_check_results (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    monitor_id VARCHAR(26) NOT NULL,
    response_status INTEGER NOT NULL,
    response_time INTEGER NOT NULL,

    CONSTRAINT fk_monitor_id FOREIGN KEY (monitor_id) REFERENCES monitors (id)
        ON DELETE CASCADE
);

CREATE TABLE monitor_request_headers (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    monitor_id VARCHAR(26) NOT NULL,
    name VARCHAR(128) NOT NULL,
    value VARCHAR(256) NOT NULL,

    CONSTRAINT fk_monitor_id FOREIGN KEY (monitor_id) REFERENCES monitors (id)
        ON DELETE CASCADE
);

CREATE TABLE regions (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE
);

CREATE TABLE monitor_regions (
    monitor_id VARCHAR(26) NOT NULL,
    region_id VARCHAR(26) NOT NULL,

    CONSTRAINT fk_monitor_id_region_id PRIMARY KEY (monitor_id, region_id),
    CONSTRAINT fk_monitor_has_regions FOREIGN KEY (monitor_id) REFERENCES monitors(id),
    CONSTRAINT fk_region_belongs_to_monitor FOREIGN KEY (region_id) REFERENCES regions(id)
);

CREATE TABLE alert_triggers (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE
);

CREATE TABLE monitor_alert_triggers (
    monitor_id VARCHAR(26) NOT NULL,
    trigger_id VARCHAR(26) NOT NULL,

    CONSTRAINT fk_monitor_id_trigger_id PRIMARY KEY (monitor_id, trigger_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE monitor_request_headers;
DROP TABLE monitor_alert_triggers;
DROP TABLE monitor_check_results;
DROP TABLE monitor_regions;
DROP TABLE monitors;
DROP TABLE alert_triggers;
DROP TABLE regions;

DROP TABLE team_members;
DROP TABLE users;
DROP TABLE teams;
DROP TABLE accounts;
DROP TYPE REQUEST_METHOD;
DROP TYPE CHECK_STATUS;
-- +goose StatementEnd
