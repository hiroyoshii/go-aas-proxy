CREATE TABLE aas (
global_id VARCHAR(255) PRIMARY KEY,
content JSONB NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL
);

CREATE TABLE submodel_proxy (
global_id VARCHAR(255) PRIMARY KEY,
short_id VARCHAR(255) NOT NULL,
semantic_id VARCHAR(255) NOT NULL,
aas_id VARCHAR(255),
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
FOREIGN KEY (aas_id) REFERENCES aas (global_id) ON DELETE CASCADE
);

CREATE INDEX submodel_proxy_idx ON submodel_proxy (short_id);