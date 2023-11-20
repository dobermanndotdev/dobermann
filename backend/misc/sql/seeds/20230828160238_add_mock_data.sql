-- +goose Up
-- +goose StatementBegin
INSERT INTO accounts (id) VALUES ('01H8YBGZF8QY158H4NP4NDXFHD');
INSERT INTO teams (id, account_id, "name")
VALUES ('01HAF76XEDP97014GMFPKXACY0', '01H8YBGZF8QY158H4NP4NDXFHD', 'Demo');

INSERT INTO members (id, account_id, first_name, last_name, primary_phone_number, secondary_phone_number, email)
VALUES ('01HAF7A5F6NMND2P86MG794D1B', '01H8YBGZF8QY158H4NP4NDXFHD', 'Firmino', 'Changani', '+351968124928', '', 'firmino.changani@gmail.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM members WHERE id = '01HAF7A5F6NMND2P86MG794D1B';
DELETE FROM teams WHERE id = '01HAF76XEDP97014GMFPKXACY0';
DELETE FROM accounts WHERE id = '01H8YBGZF8QY158H4NP4NDXFHD';
-- +goose StatementEnd
