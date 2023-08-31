{{define "custMultiCol"}}
"{{.Label}}", 
IFNULL((
SELECT JSON_GROUP_ARRAY(JSON_OBJECT(
	'value', value, 
	'id', lower({{.Table}}.id), 
	'uri', "/" || ltrim("{{.Label}}/", '#') || {{.Table}}.id))
FROM {{.Table}}
WHERE {{.Table}}.id 
IN (SELECT value
	FROM {{.JoinTable}}
	WHERE book=books.id)
), '[]')
{{end}}
