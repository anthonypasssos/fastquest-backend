-- +goose Up
-- +goose StatementBegin
ALTER TABLE question_set RENAME COLUMN creation_date TO created_at;
ALTER TABLE comment RENAME COLUMN creation_date TO created_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE question_set RENAME COLUMN created_at TO creation_date;
ALTER TABLE comment RENAME COLUMN created_at TO creation_date;
-- +goose StatementEnd
