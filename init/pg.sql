DROP TABLE IF EXISTS test_table;

DROP SEQUENCE IF EXISTS test_table_seq;

CREATE SEQUENCE test_table_seq START 1;

CREATE TABLE IF NOT EXISTS test_table(
    id INTEGER PRIMARY KEY DEFAULT nextval('test_table_seq'),
    username BYTEA,
    password BYTEA,
    email BYTEA
);
