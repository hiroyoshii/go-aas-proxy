CREATE TABLE digital_product (
    id VARCHAR(255) PRIMARY KEY,
    product_id VARCHAR(255) NOT NULL,
    serial_number VARCHAR(255),
    manufacturer_id VARCHAR(255),
    FOREIGN KEY (manufacturer_id) REFERENCES manufacturer (id) ON DELETE CASCADE
);
