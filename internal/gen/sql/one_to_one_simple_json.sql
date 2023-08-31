{{define "OneToOneSimple"}}
const {{.Label}}SimpleSql = `
(
SELECT JSON_QUOTE(IFNULL({{.Column}}, ''))
FROM {{.Table}} 
WHERE book=books.id) {{.Label}}
`
{{end}}
