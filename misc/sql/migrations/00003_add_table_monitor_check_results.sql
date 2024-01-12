-- +goose Up
-- +goose StatementBegin
CREATE TABLE monitor_check_results (
    id VARCHAR(26) PRIMARY KEY NOT NULL,
    monitor_id VARCHAR(26) NOT NULL,
    status_code SMALLINT,
    region VARCHAR(64) NOT NULL,
    response_time_in_ms SMALLINT NOT NULL,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_check_result_belongs_to_monitor FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE monitor_check_results;
-- +goose StatementEnd
