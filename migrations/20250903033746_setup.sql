-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE subject (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE question (
    id SERIAL PRIMARY KEY,
    subject_id INT REFERENCES subject(id),
    user_id INT REFERENCES users(id),
    statement TEXT NOT NULL,
    subject INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE source (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    metadata JSON
);

CREATE TABLE question_source (
    id SERIAL PRIMARY KEY,
    question_id INT REFERENCES question(id),
    source_id INT REFERENCES source(id)
);

CREATE TABLE topic (
    id SERIAL PRIMARY KEY,
    subject_id INT REFERENCES subject(id),
    name TEXT NOT NULL
);

CREATE TABLE question_topic (
    id SERIAL PRIMARY KEY,
    question_id INT REFERENCES question(id),
    topic INT REFERENCES topic(id)
);

CREATE TABLE answer (
    id SERIAL PRIMARY KEY,
    id_question INT REFERENCES question(id),
    text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL
);

CREATE TABLE question_set (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    type VARCHAR(100) NOT NULL,
    name VARCHAR(100),
    description TEXT,
    is_private BOOLEAN NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE question_set_question (
    id SERIAL PRIMARY KEY,
    question_id INT NOT NULL REFERENCES question(id),
    question_set_id INT NOT NULL REFERENCES question_set(id),
    position INT
);

CREATE TABLE user_response (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    question_id INT NOT NULL REFERENCES question(id),
    question_set_id INT NOT NULL REFERENCES question_set(id),
    answer_id INT NOT NULL REFERENCES answer(id),
    is_correct BOOLEAN NOT NULL,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comment (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    text TEXT NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comment_relationship (
    id_comment INT PRIMARY KEY REFERENCES comment(id),
    id_reference INT NOT NULL,
    type_reference TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users, subject, question, source, question_source, topic, question_topic, answer, question_set, question_set_question, user_response, comment, comment_relationship;
-- +goose StatementEnd
