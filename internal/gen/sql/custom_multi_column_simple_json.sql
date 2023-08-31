{{define "custMultiColSimple"}}
(SELECT
	CASE JSON_EXTRACT(val, "$.#{{.Label}}.is_multiple.list_to_ui") = " & "
	WHEN true
	THEN (
		SELECT JSON_QUOTE(IFNULL(GROUP_CONCAT(value, ' & '), ""))
		FROM {{.Table}}
		WHERE {{.Table}}.id IN (
			SELECT value
		FROM {{.JoinTable}}
		WHERE book=books.id)
	)
	ELSE (
		SELECT JSON_QUOTE(IFNULL(GROUP_CONCAT(value, ', '), ''))
	FROM {{.Table}}
	WHERE {{.Table}}.id IN (
		SELECT value
		FROM {{.JoinTable}}
		WHERE book=books.id)
	)
	END
FROM preferences
WHERE key = "field_metadata") "#{{.Label}}"
{{end}}
