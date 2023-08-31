{{define "custSingleCol"}}
"{{.Label}}", 
IFNULL((
SELECT JSON_QUOTE(value)
FROM {{.Table}}
WHERE book=books.id
), "")
{{end}}
