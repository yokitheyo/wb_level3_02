-- +goose Up
UPDATE urls SET is_disabled = false WHERE is_disabled IS NULL;
ALTER TABLE urls ALTER COLUMN is_disabled SET DEFAULT false;
ALTER TABLE urls ALTER COLUMN is_disabled SET NOT NULL;

-- +goose Down
ALTER TABLE urls ALTER COLUMN is_disabled DROP NOT NULL;
ALTER TABLE urls ALTER COLUMN is_disabled DROP DEFAULT;