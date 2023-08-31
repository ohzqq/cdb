{{define "OneToOneJson"}}
const {{.Label}}JsonSql = `
IFNULL((
SELECT JSON_QUOTE({{.Column}}) 
FROM {{.Table}} 
WHERE book=books.id), '') {{.Label}}
`
{{end}}
