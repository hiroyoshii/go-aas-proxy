CREATE TABLE digital_product (
    id VARCHAR(255) PRIMARY KEY,
    product_id VARCHAR(255) NOT NULL,
    serial_number VARCHAR(255),
    country_of_origin VARCHAR(10),
    construction_year INTEGER
);
