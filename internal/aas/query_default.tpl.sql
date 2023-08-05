SELECT row_to_json(aass)
FROM (
    SELECT
    	aas.short_id as idShort,
        (
        	SELECT row_to_json(nested_section)
        	FROM (
	        	SELECT
		     		global_id as id,
		     		global_id_type as idType
		        FROM aas as nested
		        WHERE nested.global_id = aas.global_id
        	) AS nested_section
        ) AS identification,
        (
        	SELECT json_agg(nested_desc)
        	FROM (
	        	SELECT
					'jp' as language,
					description as text
		        FROM aas as nested
		        WHERE nested.global_id = aas.global_id
        	) AS nested_desc
        ) AS description,
		(
        	SELECT row_to_json(nested_model)
        	FROM (
	        	SELECT
					model_type as name
		        FROM aas as nested
		        WHERE nested.global_id = aas.global_id
        	) AS nested_model
        ) AS modelType
    FROM aas
{{- if .AasID }} 
    WHERE aas.global_id = '{{ .AasID }}' 
{{- end }}
) AS aass;