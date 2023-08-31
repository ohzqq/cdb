{{define "OneToOneStr"}}
const {{.Label}}StringSql = `
(
SELECT IFNULL({{.Column}}, '')
FROM {{.Table}} 
WHERE book=books.id) {{.Label}}
`
{{end}}
