{{define "ManyToManyJson"}}
const {{.Label}}JsonSql = `
IFNULL((
SELECT 
JSON_GROUP_ARRAY(JSON_OBJECT(
'{{.Column}}', {{.Column}},
'uri', '{{.Label}}/' || lower(id),
'id', lower(id)
))
FROM {{.Table}}
WHERE {{.Table}}.id 
IN (
	SELECT {{.LinkColumn}}
	FROM {{.JoinTable}}
	WHERE book=books.id)
), '[]') {{.Label}}
`
{{end}}
