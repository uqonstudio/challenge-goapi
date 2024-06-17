INSERT INTO ms_customer (
    name, phoneNumber, address) 
VALUES(
    'Ahmad', '085678910', 'gg. dukuh I') RETURNING id;

INSERT INTO ms_employee (
    name, phoneNumber, address, email, password, department) 
VALUES(
    'Mirna', '08567889910', 'jl. pisang no A300', 'mirna@gmail.com', '589dbb5629b1baf681cfc3a7818fe5d5', 'admin') RETURNING id;

INSERT INTO ms_products(
    name, price, unit) 
VALUES('Cuci dan Strika', 7000, 'KG'),
('Laundry Spray', 50000, 'KG'),
('Laundry Tas', 25000, 'KG');
 
INSERT INTO tx_bill(
    billId, entryDate, finishDate, employee, customer, totalBill) 
VALUES('el/1/1-2024.168.22.52.00-1', '17-06-2024', '18-06-2024',1, 2, 135000);

INSERT INTO tx_billDetails(
    billId, product, productPrice, qty) 
VALUES
('el/1/1-2024.168.22.52.00-1', 1, 7000, 5),
('el/1/1-2024.168.22.52.00-1', 2, 50000,2),
('el/1/1-2024.168.22.52.00-1', 3, 25000,1);

-- DML SELECT
SELECT * FROM ms_customer WHERE id = 1;
SELECT * FROM ms_employee WHERE id = 1;
SELECT * FROM ms_products ORDER BY id ASC ;
SELECT * FROM tx_bill ORDER BY billId ASC ;
SELECT * FROM tx_billDetails ORDER BY billId ASC ;

-- DML DELETE
DELETE FROM ms_customer WHERE id = 1;
DELETE FROM ms_employee WHERE id = 1;
DELETE FROM ms_products WHERE id = 1; 


-- DML UPDATE
UPDATE ms_customer SET name = 'Ahmad' WHERE id = 1;
UPDATE ms_employee SET name = 'Mirna' WHERE id = 1;
UPDATE ms_products SET name = 'Cuci dan Strika' WHERE id = 1;


