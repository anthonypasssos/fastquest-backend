-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, name, email, password_hash) VALUES
(1, 'Beatrice Yure', 'bea.yure@gmail.com', 'TKATRTKVAEAD'),
(2, 'Davi Pingado', 'pingado.davi.tsun@gmail.com', 'YSETHBS'),
(3, 'Guilherme Ashevale', 'guilherme.ashevale@gmail.com', 'VWARAKFRYW'),
(4, 'Beatriz Dere', 'beatriz.dere@gmail.com', 'VEABIYRYUIB'),
(5, 'FlashAi', 'gemini.flashquest@gmail.com', 'HFHEOHDODHWOQIHD');

INSERT INTO subject (id, name) VALUES
(1, 'Direito Constitucional'),
(2, 'Direito Penal'),
(3, 'Direito Civil'),
(4, 'Direito Administrativo'),
(5, 'Direito do Trabalho'),
(6, 'Direito Tributário'),
(7, 'AI');

INSERT INTO question (id, statement, subject_id, user_id) VALUES
(1, 'Determinada sociedade de advogados deseja se associar a advogados que não a integram para prestação de serviços e participação nos resultados. Segundo a legislação aplicável à formalização desse vínculo jurídico, assinale a opção que indica, corretamente, a conclusão dos administradores da sociedade de advogados.', 3, 1),
(2, 'Sebastião, advogado, celebrou contrato de mandato com o cliente Amir, para representá-lo extrajudicialmente, tendo realizado diligências em prol da resolução do imbróglio. Desde a celebração do mandato, passaram-se mais de 20 (vinte) anos, mas as atividades para as quais Amir contratou Sebastião, por sua própria natureza, se protraíram no tempo, sendo ainda necessárias a Amir. Sobre a hipótese apresentada, assinale a afirmativa correta.', 3, 1);

INSERT INTO answer (text, is_correct, id_question) VALUES
('O contrato de associação não pode ser pactuado em caráter geral, devendo restringir-se a causas ou trabalhos específicos, sob pena de se configurarem os requisitos legais de vínculo empregatício.', TRUE, 1),
('O contrato de associação deverá ser registrado no Conselho Seccional da OAB em cuja base territorial tiver sede a sociedade de advogados.', FALSE, 1),
('O contrato de associação poderá atribuir a totalidade dos riscos à sociedade de advogados, mas não exclusivamente a um advogado sócio ou associado.', FALSE, 1),
('O advogado não pode, simultaneamente, celebrar contrato de associação com mais de uma sociedade de advogados com sede ou filial na mesma área territorial do respectivo Conselho Seccional.', FALSE, 1),
('O mandato extinguiu-se pelo decurso do tempo, salvo se previsto prazo diverso no respectivo instrumento.', TRUE, 2),
('O mandato extinguiu-se pelo decurso do tempo, sendo vedada a previsão de prazo diverso no respectivo instrumento.', FALSE, 2),
('O mandato não se extinguiu pelo decurso do tempo, salvo se foi consignado prazo no respectivo instrumento.', FALSE, 2),
('O mandato não se extinguiu pelo decurso do tempo, sendo vedada a estipulação de prazo no respectivo instrumento.', FALSE, 2);

INSERT INTO question (id, statement, subject_id, user_id) VALUES
(3, 'Monique, advogada regularmente inscrita nos quadros da OAB, é investigada em inquérito policial por supostos crimes praticados por motivo ligado ao exercício da advocacia, tendo sido presa em flagrante, por crime da mesma espécie, em seu escritório, enquanto atendia a uma de suas clientes. Considerando as disposições do Estatuto da Advocacia, é correto afirmar que', 3, 1),
(4, 'Pedro, contador com vasta experiência e sólida carreira, decide fazer uma segunda graduação, tornando-se bacharel em Direito. Depois da aprovação no Exame de Ordem Unificado e da inscrição nos quadros da Ordem dos Advogados do Brasil, Pedro pretende continuar prestando serviços contábeis, sem prejuízo do exercício concomitante da nova atividade. Acerca da intenção de Pedro, bem como dos limites ético-normativos para a publicidade profissional da sua nova atividade, assinale a afirmativa correta.', 3, 1);

INSERT INTO answer (text, is_correct, id_question) VALUES
('Monique tem direito à presença de representante da OAB para lavratura do auto de prisão em flagrante, visto que se trata de suposto crime por motivo ligado ao exercício da advocacia, sob pena de nulidade.', TRUE, 3),
('Não há qualquer direito ou prerrogativa conferida pela legislação no caso em tela, devendo Monique receber tratamento idêntico ao dado a outros indivíduos não advogados, em razão do princípio da igualdade.', FALSE, 3),
('A presença de representante da OAB no momento da lavratura do auto de prisão em flagrante será devida ainda que não se trate de motivo ligado ao exercício da advocacia, visto que se cuida de direito conferido ao advogado em todo e qualquer crime por ele cometido.', FALSE, 3),
('O representante da OAB para acompanhar a lavratura do auto de prisão em flagrante, pode ser substituído por representante da Defensoria Pública, visto que ambos podem figurar como defensores.', FALSE, 3),
('Pedro não poderá exercer de modo concomitante as atividades de contador e advogado, pois, de acordo com o Estatuto da Advocacia e da OAB, a prestação de serviços contábeis é incompatível com o exercício simultâneo da advocacia.', FALSE, 4),
('Não há óbice ético para o duplo exercício das atividades de contador e advogado, podendo Pedro se valer da divulgação conjunta dos serviços oferecidos, desde que não seja por meio de inscrições em muros, paredes, veículos, elevadores ou em qualquer espaço público.', TRUE, 4),
('Embora não haja incompatibilidade para o exercício concomitante das duas atividades, não será permitido a Pedro divulgar sua nova profissão de modo conjunto com a de contador.', FALSE, 4),
('Pedro poderá fazer uso de mala direta, distribuição de panfletos ou formas assemelhadas de publicidade, visando a captação de clientela para a sua nova atividade, mas não poderá mencionar, nessa publicidade, os serviços de contabilidade.', FALSE, 4);

INSERT INTO question (id, statement, subject_id, user_id) VALUES
(5, 'Formalizou-se, no Tribunal Regional Eleitoral do Estado Alfa, a vacância de um dos cargos de juiz eleitoral, reservado constitucionalmente à classe de advogados. De igual modo, no Tribunal Regional Federal da Enésima Região, sediado na capital do mesmo Estado Alfa, com jurisdição nos Estados Alfa, Beta e Gama, foi também formalizada a vacância de um cargo de juiz federal do Tribunal Regional Federal, destinado à advocacia nos termos da Constituição Federal (quinto constitucional). Nesse hipotético cenário, que demandará a produção de duas listas de membros da advocacia para o futuro preenchimento dos cargos, assinale a afirmativa que descreve corretamente as competências dos órgãos da OAB.', 3, 1),
(6, 'Valmir, bacharel em Direito, aprovado no Exame da Ordem dos Advogados do Brasil, ocupa o cargo público de agente de Polícia Civil do Estado Alfa. Movido por sentimento altruísta, Valmir requer sua inscrição na OAB, pois pretende, nos momentos de folga da atividade policial, exercer a advocacia de forma gratuita, eventual e voluntária, em favor de instituições sociais sem fins econômicos que não disponham de recursos para a contratação de profissional. À luz dessas informações, e considerada a legislação vigente, assinale a afirmativa correta.', 3, 1);

INSERT INTO answer (text, is_correct, id_question) VALUES
('A lista para o preenchimento do cargo no TRE do Estado Alfa ficará sob a incumbência do Conselho Seccional da OAB do respectivo Estado, competindo ao Conselho Federal da OAB elaborar a lista para o preenchimento do cargo no TRF da Enésima Região.', TRUE, 5),
('A lista para o preenchimento do cargo no TRE do Estado Alfa ficará sob a incumbência do Conselho Seccional da OAB do respectivo Estado, competindo aos Conselhos Seccionais da OAB dos Estados Alfa, Beta e Gama a elaboração conjunta da lista para o preenchimento do cargo no TRF da Enésima Região.', FALSE, 5),
('Uma vez que tanto a Justiça Eleitoral quanto a Justiça Federal pertencem ao Poder Judiciário da União, competirá ao Conselho Federal da OAB a elaboração das duas listas.', FALSE, 5),
('Uma vez que tanto o TRE do Estado Alfa quanto a sede do TRF da Enésima Região estão situados no Estado Alfa, competirá ao Conselho Seccional da OAB desse Estado a elaboração das duas listas.', FALSE, 5),
('Valmir poderá exercer regularmente a advocacia, inclusive pro bono.', FALSE, 6),
('Valmir não poderá exercer a advocacia remunerada, pois ocupa cargo incompatível, mas poderá exercer a advocacia pro bono.', TRUE, 6),
('Valmir não poderá exercer a advocacia, mesmo pro bono, uma vez que o cargo público que ocupa atrai o regime da incompatibilidade.', FALSE, 6),
('A condição de servidor público atrai o regime do impedimento, razão pela qual Valmir não poderá exercer a advocacia contra a Fazenda Pública que o remunera. Observado esse impedimento, não haverá óbice para o exercício da advocacia pro bono.', FALSE, 6);

SELECT setval('question_id_seq', (SELECT MAX(id) FROM question));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE answer, question, subject, users RESTART IDENTITY CASCADE;
-- +goose StatementEnd
