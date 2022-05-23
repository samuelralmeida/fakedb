CREATE DATABASE IF NOT EXISTS prod_trem;
CREATE TABLE IF NOT EXISTS  prod_trem.client (
	id int,
	name varchar(255),
	cpf varchar(255),
	address varchar(255)
);
INSERT INTO
	prod_trem.client (id, name, cpf, address)
VALUES
	(1, "Arthur Dent", "49735350017", "Rua Urano, 1554"),
	(2, "Ford Prefect", "97100265002", "Avenida Jupiter, 42545" ),
	(3, "Zaphod Beeblebrox", "86290863002", "Travessa Marte, 518")
ON DUPLICATE KEY UPDATE id=id;

CREATE DATABASE IF NOT EXISTS prod_coisa;
CREATE TABLE IF NOT EXISTS prod_coisa.leads (
	id int,
	email varchar(255)
);
INSERT INTO
	prod_coisa.leads (id, email)
VALUES
	(1, "marvin@hitchhiker.galaxy"),
	(2, "trillian@hitchhiker.galaxy"),
	(3, "vogon@hitchhiker.galaxy")
ON DUPLICATE KEY UPDATE id=id;