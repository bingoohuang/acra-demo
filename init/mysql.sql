CREATE DATABASE IF NOT EXISTS test;

use test;

DROP TABLE IF EXISTS test_table;

CREATE TABLE IF NOT EXISTS test_table(
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARBINARY(2000),
    password VARBINARY(2000),
    email VARBINARY(2000)
);
