-- +goose Up
-- +goose StatementBegin
TRUNCATE TABLE 
    users,
    subject,
    question,
    answer
RESTART IDENTITY CASCADE;

INSERT INTO users (id, name, email, password_hash) VALUES
(1, 'Sistema', 'system.flashquest@gmail.com', 'v4QsVeX2ztDbhqpXepQDuBlq3zgr2FA0FKMrTFJ105tdd1d9yd'),
(2, 'FlashAi', 'ai.flashquest@gmail.com', 'k7ZIKXS2IaaOKM5U0Wrx9sir7TJH2epiv4kePQUFTQdRegXiVx');

INSERT INTO subject (id, name) VALUES
(1, 'Ética Profissional'),
(2, 'Filosofia do Direito'),
(3, 'Direito Constitucional'),
(4, 'Direitos Humanos'),
(5, 'Direito Eleitoral'),
(6, 'Direito Internacional'),
(7, 'Direito Financeiro'),
(8, 'Direito Tributário'),
(9, 'Direito Administrativo'),
(10, 'Direito Ambiental'),
(11, 'Direito Civil'),
(12, 'ECA'),
(13, 'Direito do Consumidor'),
(14, 'Direito Empresarial'),
(15, 'Processo Civil'),
(16, 'Direito Penal'),
(17, 'Processo Penal'),
(18, 'Direito Previdenciário'),
(19, 'Direito do Trabalho'),
(20, 'Processo do Trabalho');

INSERT INTO source (id, name, type) VALUES (1, 'OAB', 'exam');

SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('subject_id_seq', (SELECT MAX(id) FROM subject));
SELECT setval('source_id_seq', (SELECT MAX(id) FROM source));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE 
    users,
    subject,
    question,
    answer,
    source,
    source_exam_instance
RESTART IDENTITY CASCADE;
-- +goose StatementEnd
