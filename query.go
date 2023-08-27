package cdb

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Query struct {
	query sq.SelectBuilder
	sort  string
	limit int
}

func NewQuery(cols []string) *Query {
	q := &Query{
		sort: "timestamp",
	}

	q.query = sq.Select(strings.Join(cols, ",\n")).
		From("books")

	return q
}

func (q *Query) GetByID(ids ...any) *Query {
	q.query = q.query.Where(sq.Eq{"id": ids})
	return q
}

func (q *Query) Sort(v string) *Query {
	db := strings.Split(v, ":")
	sort := db[0]
	if order := db[1]; order == "desc" {
		sort += " DESC\n"
	}
	q.query = q.query.OrderBy(sort)
	return q
}

func (q *Query) Limit(v int) *Query {
	q.limit = v
	q.query = q.query.Limit(uint64(v))
	return q
}

func (q *Query) Page(v int) *Query {
	if q.limit == 0 {
		q.Limit(20)
	}
	q.query = q.query.Offset(calculateOffset(v, q.limit))
	return q
}

func (q *Query) Filter(w string) *Query {
	q.query = q.query.Where(w)
	return q
}

func (q *Query) ToSql() (string, []any, error) {
	return q.query.ToSql()
}

func AddedInThePastDays(d int) string {
	if d == 0 {
		d = 7
	}
	return fmt.Sprintf("last_modified > DATE('now', '-%d day')", d)
}
