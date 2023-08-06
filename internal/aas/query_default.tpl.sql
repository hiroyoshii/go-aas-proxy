SELECT 
	content
FROM aas
{{- if .AasID }} 
WHERE aas.global_id = '{{ .AasID }}' 
{{- end }}
;