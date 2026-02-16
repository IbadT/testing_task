SELECT 'CREATE DATABASE testingtask'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'testingtask')\gexec

CREATE EXTENSION IF NOT EXISTS pgcrypto;

\echo '✅ init.sql completed: database and extensions created'