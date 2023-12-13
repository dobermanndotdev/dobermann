-- +goose Up
-- +goose StatementBegin
CREATE TABLE monitor_check_results (
    id SERIAL PRIMARY KEY NOT NULL,
    monitor_id VARCHAR(26) NOT NULL,
    status_code SMALLINT NOT NULL,
    region VARCHAR(64) NOT NULL,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    response_time_in_ms SMALLINT NOT NULL,

    CONSTRAINT fk_check_result_belongs_to_monitor FOREIGN KEY (monitor_id) REFERENCES monitors(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE monitor_check_results;
-- +goose StatementEnd
