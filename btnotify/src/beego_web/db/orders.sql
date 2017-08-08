CREATE DATABASE orders;

create table orders (
	`sender` VARCHAR(128) NOT NULL,
	`hash` VARCHAR(128) NOT NULL,
	`value` INT(10) NOT NULL DEFAULT 0,
	`verify` INT(10) NOT NULL DEFAULT 0,
	INDEX `sender` (`sender`)
);