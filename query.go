package cdb

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// Query represents a database query.
type Query struct {
	query sq.SelectBuilder
	sort  string
	limit int
}

func newQuery(cols []string) *Query {
	q := &Query{
		sort: "timestamp",
	}

	q.query = sq.Select(strings.Join(cols, ",\n")).
		From("books")

	return q
}

// NewQuery initializes a query builder.
func NewQuery(cols []string) *Query {
	q := &Query{
		sort: "timestamp",
	}

	q.query = sq.Select("JSON_OBJECT( \n" + strings.Join(cols, ",\n") + ")").
		From("books")

	return q
}

// GetByID retrieves records by id.
func (q *Query) GetByID(ids ...any) *Query {
	q.query = q.query.Where(sq.Eq{"id": ids})
	return q
}

// Sort sets Sort By. Takes a string argument in the form of 'sort:order'.
func (q *Query) Sort(v string) *Query {
	db := strings.Split(v, ":")
	sort := db[0]
	if order := db[1]; order == "desc" {
		sort += " DESC\n"
	}
	q.query = q.query.OrderBy(sort)
	return q
}

// Limit sets the number of records returned.
func (q *Query) Limit(v int) *Query {
	q.limit = v
	q.query = q.query.Limit(uint64(v))
	return q
}

// Page is used with the limit to determine OFFSET.
func (q *Query) Page(v int) *Query {
	if q.limit == 0 {
		q.Limit(20)
	}
	q.query = q.query.Offset(calculateOffset(v, q.limit))
	return q
}

// Filter takes an arbitrary WHERE expression.
func (q *Query) Filter(w string) *Query {
	q.query = q.query.Where(w)
	return q
}

// ToSql satisfies squirrel.Sqlizer to build the sql expression with paramaters.
func (q *Query) ToSql() (string, []any, error) {
	return q.query.ToSql()
}

// AddedInThePastDays retrieves the past n days of records.
func AddedInThePastDays(d int) string {
	if d == 0 {
		d = 7
	}
	return fmt.Sprintf("last_modified > DATE('now', '-%d day')", d)
}
