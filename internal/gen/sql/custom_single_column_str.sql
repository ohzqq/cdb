{{define "custSingleColStr"}}
(SELECT
	CASE JSON_EXTRACT(val, "$.{{.Label}}.is_multiple.list_to_ui") = " & "
	WHEN true
	THEN (
		SELECT IFNULL(GROUP_CONCAT(value, ' & '), "")
		FROM {{.Table}}
		WHERE book=books.id
		) 
	ELSE (
		SELECT IFNULL(GROUP_CONCAT(value, ', '), '')
		FROM {{.Table}}
		WHERE book=books.id
		)
	END
FROM preferences
WHERE key = "field_metadata") "{{.Label}}"
{{end}}
