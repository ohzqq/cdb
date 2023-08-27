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

var Preferences = []string{
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
