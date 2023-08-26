package cdb

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Search struct {
	query sq.SelectBuilder
	db    *DB
	sort  string
	limit int
}

func (db *Search) GetByID(ids ...any) *Search {
	db.query = db.query.Where(sq.Eq{"id": ids})
	return db
}

func (q *Search) Sort(v string) *Search {
	db := strings.Split(v, ":")
	sort := db[0]
	if order := db[1]; order == "desc" {
		sort += " DESC\n"
	}
	q.query = q.query.OrderBy(sort)
	return q
}

func (db *Search) Limit(v int) *Search {
	db.limit = v
	db.query = db.query.Limit(uint64(v))
	return db
}

func (db *Search) Page(v int) *Search {
	if db.limit == 0 {
		db.Limit(20)
	}
	db.query = db.query.Offset(calculateOffset(v, db.limit))
	return db
}

func (db *Search) Filter(w string) *Search {
	db.query = db.query.Where(w)
	return db
}

func (s *Search) Results() ([]*Row, error) {
	return s.db.execute(toSql(s.query))
}

func AddedInThePastDays(d int) string {
	if d == 0 {
		d = 7
	}
	return fmt.Sprintf("last_modified > DATE('now', '-%d day')", d)
}
