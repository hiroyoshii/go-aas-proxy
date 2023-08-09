CREATE TABLE manufacturer (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(30),
    gl_number VARCHAR(255),
    supplier_name VARCHAR(255)
);

CREATE TABLE product (
    id VARCHAR(255) PRIMARY KEY,
    serial_number VARCHAR(255),
    batch_number VARCHAR(255),
    manufacturer_id VARCHAR(255),
    FOREIGN KEY (manufacturer_id) REFERENCES manufacturer (id) ON DELETE CASCADE
);

CREATE TABLE contact_info (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255),
    url  VARCHAR(255),
    phone_number  VARCHAR(255),
    fax  VARCHAR(255),
    manufacturer_id VARCHAR(255),
    FOREIGN KEY (manufacturer_id) REFERENCES manufacturer (id) ON DELETE CASCADE
);
