INSERT INTO ms_customer (
    name, phoneNumber, address) 
VALUES(
    'Ahmad', '085678910', 'gg. dukuh I') RETURNING id;

INSERT INTO ms_employee (
    name, phoneNumber, address, email, password, department) 
VALUES(
    'Mirna', '08567889910', 'jl. pisang no A300', 'mirna@gmail.com', '589dbb5629b1baf681cfc3a7818fe5d5', 'admin') RETURNING id;

INSERT INTO ms_products(
    service_name, unit, price) 
VALUES('Cuci + Strika', 'KG', 7000) RETURNING id;
 
INSERT INTO tx_bill(
    billId, entryDate, finishDate, employee, customer, totalBill) 
VALUES('el/1/1-2024.168.22.52.00-1', '17-06-2024', '18-06-2024',1, 2, 135000);

INSERT INTO tx_billDetails(
    billId, product, productPrice, qty) 
VALUES
('el/1/1-2024.168.22.52.00-1', 1, 7000, 5),
('el/1/1-2024.168.22.52.00-1', 2, 50000,2),
('el/1/1-2024.168.22.52.00-1', 3, 25000,1);
