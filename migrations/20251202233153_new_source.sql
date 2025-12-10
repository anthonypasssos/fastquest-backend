-- +goose Up
-- +goose StatementBegin
ALTER TABLE source DROP COLUMN metadata;
ALTER TABLE source ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

CREATE TABLE source_exam_instance(
    id SERIAL PRIMARY KEY,
    source_id INT NOT NULL REFERENCES source(id),
    edition INT NOT NULL,
    phase INT NOT NULL,
    year INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE question
ADD COLUMN source_exam_instance_id BIGINT REFERENCES source_exam_instance(id);
ALTER TABLE question DROP COLUMN subject;

DROP TABLE IF EXISTS question_sources;
DROP TABLE question_source;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE question DROP COLUMN source_exam_instance_id;
DROP TABLE IF EXISTS source_exam_instance;

ALTER TABLE source DROP COLUMN IF EXISTS created_at;
ALTER TABLE source ADD COLUMN metadata JSON;

CREATE TABLE question_source (
    id SERIAL PRIMARY KEY,
    question_id INT REFERENCES question(id),
    source_id INT REFERENCES source(id)
);
-- +goose StatementEnd
