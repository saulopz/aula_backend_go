-- Estrutura do banco da usersapi.
--
-- Para criar o banco e a tabela:
--   createdb usersapi
--   psql -d usersapi -f schema.sql
--
-- IF NOT EXISTS: pode executar de novo sem erro.

CREATE TABLE IF NOT EXISTS users (
    -- O próprio banco gera o ID, somando 1 a cada registro.
    -- ALWAYS: o banco recusa um INSERT que tente escolher o ID.
    id    INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name  VARCHAR(50)  NOT NULL,
    email VARCHAR(100) NOT NULL
);
