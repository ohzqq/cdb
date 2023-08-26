package cdb

import (
	sq "github.com/Masterminds/squirrel"
)

type Search struct {
	query sq.SelectBuilder
	sort  string
	limit int
}
