{{define "ManyToManySimple"}}
const {{.Label}}SimpleSql = `
(
SELECT 
{{if .IsNames -}}
JSON_QUOTE(IFNULL(GROUP_CONCAT({{.Column}}, ' & '), ''))
{{- else -}}
JSON_QUOTE(IFNULL(GROUP_CONCAT({{.Column}}, ', '), ''))
{{- end}}
FROM {{.Table}}
WHERE {{.Table}}.id 
IN (
	SELECT {{.LinkColumn}}
	FROM {{.JoinTable}}
	WHERE book=books.id)
) {{.Label}}
`
{{end}}
