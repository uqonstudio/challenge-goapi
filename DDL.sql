CREATE DATABASE enigmalaundry;

CREATE TABLE ms_customer(
id SERIAL,
name VARCHAR(100) NOT NULL,
phoneNumber VARCHAR(100) NOT NULL,
address VARCHAR(100),
PRIMARY KEY(id));

CREATE TABLE ms_employee(
id SERIAL,
name VARCHAR(100) NOT NULL,
password VARCHAR,
email VARCHAR(100) NOT NULL,
phoneNumber VARCHAR(100) NOT NULL,
department VARCHAR(100) NOT NULL,
PRIMARY KEY(id));

CREATE TABLE ms_products(
id SERIAL,
name VARCHAR(100) NOT NULL,
unit VARCHAR(100) NOT NULL,
price INT NOT NULL,
PRIMARY KEY(id));

CREATE TABLE tx_bill(
id SERIAL,
billId VARCHAR(100) NOT NULL,
entryDate DATE NOT NULL DEFAULT CURRENT_DATE,
finishDate DATE NOT NULL,
employee INT NOT NULL,
customer INT NOT NULL,
totalBill INT, 
PRIMARY KEY(billId),
FOREIGN KEY(employee) REFERENCES ms_employee(id),
FOREIGN KEY(customer) REFERENCES ms_customer(id));

CREATE TABLE tx_billDetails(
id SERIAL,
billId VARCHAR NOT NULL,
product INT NOT NULL,
qty INT NOT NULL, 
price INT NOT NULL,
sub_total INT NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(billId) REFERENCES tx_bill(billId),
FOREIGN KEY(product) REFERENCES ms_products(id));
