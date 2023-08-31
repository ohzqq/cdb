{{define "ManyToManyStr"}}
const {{.Label}}StringSql = `
(
SELECT 
{{if .IsNames -}}
IFNULL(GROUP_CONCAT({{.Column}}, ' & '), '')
{{- else -}}
	{{if eq .Label "rating"}}
	IFNULL(GROUP_CONCAT(lower({{.Column}}), ', '), '')
	{{else}}
	IFNULL(GROUP_CONCAT({{.Column}}, ', '), '')
	{{end}}
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
