CREATE TABLE aas (
global_id VARCHAR(255) PRIMARY KEY,
global_id_type VARCHAR(255) NOT NULL,
short_id VARCHAR(255) NOT NULL UNIQUE,
description VARCHAR(255),
model_type VARCHAR(10) NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL
);

CREATE TABLE submodel_proxy (
short_id VARCHAR(255) PRIMARY KEY,
semantic_id VARCHAR(255) NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
aas_id VARCHAR(255),
FOREIGN KEY (aas_id) REFERENCES aas (global_id) ON DELETE CASCADE
);
