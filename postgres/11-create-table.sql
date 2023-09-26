\c persons

CREATE TABLE IF NOT EXISTS tb_persons(
    id serial PRIMARY KEY,
    name varchar(128),
    age int,
    address varchar(256),
    work varchar(256)
);