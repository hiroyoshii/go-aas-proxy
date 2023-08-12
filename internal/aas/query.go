package aas

var (
	submodelRefSql = `
SELECT
 short_id,
 semantic_id 
FROM
 submodel_proxy
JOIN aas ON aas.global_id = submodel_proxy.aas_id
WHERE
 true
 {{- if ne .AasID "" } }}
 AND aas.global_id = '{{ .AasID }}'
 {{- end }}
 {{- if ne .SubmodelIDShort "" } }}
 AND submodel_proxy.short_id = '{{ .SubmodelIDShort }}'
 {{- end }}
;
`
	aasSql = `
SELECT 
	content
FROM aas
{{- if .AasID }} 
WHERE aas.global_id = '{{ .AasID }}' 
{{- end }}
;
`
	aasInitSql = `
CREATE TABLE IF NOT EXISTS aas (
	global_id VARCHAR(255) PRIMARY KEY,
	content JSONB NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
	);
	
CREATE TABLE IF NOT EXISTS submodel_proxy (
global_id VARCHAR(255) PRIMARY KEY,
short_id VARCHAR(255) NOT NULL,
semantic_id VARCHAR(255) NOT NULL,
aas_id VARCHAR(255),
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
FOREIGN KEY (aas_id) REFERENCES aas (global_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS submodel_proxy_idx ON submodel_proxy (short_id);
`
)
