package cdb

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Query struct {
	query sq.SelectBuilder
	db    *DB
	sort  string
	limit int
}

func Search(l string) *Query {
	return GetLib(l).NewQuery()
}

func (db *Query) GetByID(ids ...any) *Query {
	db.query = db.query.Where(sq.Eq{"id": ids})
	return db
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

func (db *Query) Limit(v int) *Query {
	db.limit = v
	db.query = db.query.Limit(uint64(v))
	return db
}

func (db *Query) Page(v int) *Query {
	if db.limit == 0 {
		db.Limit(20)
	}
	db.query = db.query.Offset(calculateOffset(v, db.limit))
	return db
}

func (db *Query) Filter(w string) *Query {
	db.query = db.query.Where(w)
	return db
}

func (s *Query) Results() ([]*Book, error) {
	return s.db.execute(toSql(s.query))
}

func AddedInThePastDays(d int) string {
	if d == 0 {
		d = 7
	}
	return fmt.Sprintf("last_modified > DATE('now', '-%d day')", d)
}
