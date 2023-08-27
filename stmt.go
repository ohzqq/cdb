package cdb

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func booksColumn(m Model) (string, []any) {
	return fmt.Sprintf(`IFNULL(%s, '') %s`, m.Column, m.Label), []any{}
}

func groupConcat(m Model) string {
	sep := defaultSep
	if m.Label == Authors || m.Label == Narrators {
		sep = namesSep
	}
	return fmt.Sprintf("IFNULL(GROUP_CONCAT(%s, '%s'), '')", m.Column, sep)
}

func manyToOne(m Model) (string, []any) {
	sel := sq.Select(groupConcat(m)).
		Prefix("(").
		From(m.Table).
		Where("book=books.id").
		Suffix(fmt.Sprint(") '", m.Label, "'"))

	return toSql(sel)
}

func manyToMany(m Model) (string, []any) {
	col := sq.Select(groupConcat(m)).
		From(m.Table).
		Where(fmt.Sprint(m.Table, ".id"))

	join := sq.Select(m.LinkColumn).
		From(m.JoinTable).
		Where("book IN (books.id)")

	cs, _ := toSql(col)
	js, _ := toSql(join)

	stmt := fmt.Sprintf("(%s IN (%s)) '%s'", cs, js, m.Label)

	return stmt, []any{}
}

func toSql(sel sq.SelectBuilder) (string, []any) {
	stmt, args, err := sel.ToSql()
	if err != nil {
		fmt.Errorf("%v\n", err)
	}
	return stmt, args
}

func coverStmt(lib string) (string, []any) {
	sel := sq.Case("has_cover").
		When("1", `IFNULL('books/' || books.id || '/cover.jpg', '')`).
		Else("''")

	stmt, args, _ := sel.ToSql()
	return stmt + " cover", args
}

const (
	namesSep   = " & "
	defaultSep = ", "
)

func calculateOffset(page, limit int) uint64 {
	cur := page - 1
	if cur < 0 {
		cur = 0
	}
	return uint64(cur * limit)
}

func GetPreferences(col ...string) (string, []any) {
	if len(col) == 0 {
		col = PreferenceColumns
	}

	args := make([]any, len(col))
	for i, c := range col {
		args[i] = c
		col[i] = "JSON(val) " + c
	}

	sel := sq.Select(col...).
		From("preferences").
		Where(sq.Eq{"key": args})

	return toSql(sel)
}

var PreferenceColumns = []string{
	SavedSearches,
	DisplayFields,
	HiddenCategories,
	FieldMetadata,
}

const (
	SavedSearches    = "saved_searches"
	DisplayFields    = "book_display_fields"
	HiddenCategories = "tag_browser_hidden_categories"
	FieldMetadata    = "field_metadata"
)

const customColumnsFieldMetaSql = `
SELECT
JSON_GROUP_OBJECT("#" || label,
(SELECT
JSON_OBJECT(
'category_sort', IFNULL(JSON_EXTRACT(val, "$." || "#" || label || ".category_sort"), ''),
'column', IFNULL(JSON_EXTRACT(val, "$." || "#" || label || ".column"), 'uuid'),
'is_category',
CASE JSON_EXTRACT(val, "$." || "#" || label || ".is_category")
WHEN true then JSON("true")
ELSE JSON("false")
END,
'is_custom', 
CASE JSON_EXTRACT(val, "$." || "#" || label || ".is_custom")
WHEN true then JSON("true")
ELSE JSON("false")
END,
'is_editable', 
CASE JSON_EXTRACT(val, "$." || "#" || label || ".is_editable")
WHEN true then JSON("true")
ELSE JSON("false")
END,
'is_multiple', 
CASE JSON_EXTRACT(val, "$." || "#" || label || ".is_multiple")
WHEN '{}' THEN JSON("false")
ELSE JSON("true")
END,
'is_names',
CASE JSON_EXTRACT(val, "$." || "#" || label || ".is_multiple.ui_to_list")
WHEN "&" THEN JSON("true")
ELSE JSON("false")
END,
'join_table', IFNULL(
(SELECT 
'books_' || JSON_EXTRACT(val, "$." || "#" || label || ".table") || '_link'
WHERE JSON_EXTRACT(val, "$." || "#" || label || ".table") IS NOT NULL
AND JSON_EXTRACT(val, "$." || "#" || label || ".is_category") = true
), ''),
'label', "#" || JSON_EXTRACT(val, "$." || "#" || label || ".label"),
'link_column', IFNULL(JSON_EXTRACT(val, "$." || "#" || label || ".link_column"), ''),
'name', IFNULL(lower(JSON_EXTRACT(val, "$." || "#" || label || ".name")), label),
'table', IFNULL(JSON_EXTRACT(val, "$." || "#" || label || ".table"), '')
)
FROM preferences 
WHERE key = 'field_metadata')) fieldMeta
FROM custom_columns
`
