package cdb

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqliteOpts   = "?cache=shared&mode=ro"
	sqlitePrefix = `file:`
	metaDB       = `metadata.db`
)

type DB struct {
	Name       string
	Models     Models
	db         *sqlx.DB
	mtx        sync.Mutex
	printQuery bool
	sort       string
	limit      int
	query      sq.SelectBuilder
}

type Row struct {
	Title        string  `db:"title" yaml:"title"`
	Authors      string  `db:"authors" yaml:"authors,omitempty"`
	Narrators    string  `db:"narrators" yaml:"#narrators,omitempty"`
	Series       string  `db:"series" yaml:"series,omitempty"`
	SeriesIndex  float32 `db:"series_index" yaml:"series_index,omitempty"`
	Tags         string  `db:"tags" yaml:"tags,omitempty"`
	Pubdate      string  `db:"pubdate" yaml:"pubdate,omitempty"`
	Timestamp    string  `db:"timestamp" yaml:"timestamp,omitempty"`
	Duration     string  `db:"duration" yaml:"#duration,omitempty"`
	Comments     string  `db:"comments" yaml:"comments,omitempty"`
	Rating       string  `db:"rating" yaml:"rating,omitempty"`
	Publisher    string  `db:"publisher" yaml:"publisher,omitempty"`
	Languages    string  `db:"languages" yaml:"languages,omitempty"`
	Cover        string  `db:"cover" yaml:"-"`
	Formats      string  `db:"formats" yaml:"-"`
	Identifiers  string  `db:"identifiers" yaml:"identifiers,omitempty"`
	LastModified string  `db:"last_modified" yaml:"last_modified,omitempty"`
	ID           int     `db:"id" yaml:"-"`
	AuthorSort   string  `db:"author_sort" yaml:"author_sort,omitempty"`
	Sort         string  `db:"sort" yaml:"sort,omitempty"`
	Path         string  `db:"path" yaml:"-"`
	UUID         string  `db:"uuid,omitempty" yaml:"-"`
}

func Configure(name, path string, audiobooks bool) (*DB, error) {
	db := &DB{
		Name:   name,
		Models: modelMeta,
		sort:   "timestamp",
	}

	if ok := FileExist(path); !ok {
		return db, ErrFileNotExist(path)
	}

	url := sqlitePrefix + path + sqliteOpts
	database, err := sqlx.Open("sqlite3", url)
	if err != nil {
		return db, fmt.Errorf("database connection %v failed\n", err)
	}

	db.db = database

	if audiobooks {
		db.Models = AudiobookModels()
		err := db.getAudiobookColumns()
		return db, err
	}

	return db, nil
}

func (db DB) IsConnected() bool {
	return db.db != nil
}

func (db *DB) NewSearch() *DB {
	var cols []string
	for _, m := range db.Models {
		cols = append(cols, m.ToSql())
	}

	db.query = sq.Select(strings.Join(cols, ",\n")).
		From("books")

	return db
}

func (db *DB) GetByID(ids ...any) *DB {
	db.query = db.query.Where(sq.Eq{"id": ids})
	return db
}

func (q *DB) Sort(v string) *DB {
	db := strings.Split(v, ":")
	sort := db[0]
	if order := db[1]; order == "desc" {
		sort += " DESC\n"
	}
	q.query = q.query.OrderBy(sort)
	return q
}

func (db *DB) Limit(v int) *DB {
	db.limit = v
	db.query = db.query.Limit(uint64(v))
	return db
}

func (db *DB) Page(v int) *DB {
	if db.limit == 0 {
		db.Limit(20)
	}
	db.query = db.query.Offset(calculateOffset(v, db.limit))
	return db
}

func (db *DB) Filter(w string) *DB {
	db.query = db.query.Where(w)
	return db
}

func AddedInThePastDays(d int) string {
	if d == 0 {
		d = 7
	}
	return fmt.Sprintf("last_modified > DATE('now', '-%d day')", d)
}

func (db *DB) Results() ([]*Book, error) {
	return db.execute(toSql(db.query))
}

func (db *DB) execute(stmt string, args []any) ([]*Book, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	var books []*Book

	if db.printQuery {
		fmt.Println(stmt)
	}

	rows, err := db.db.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		return books, fmt.Errorf("error %v\n", err)
	}
	defer rows.Close()
	db.db.Unsafe()

	for rows.Next() {
		b := &Row{}
		err := rows.StructScan(b)
		if err != nil {
			return books, err
		}
		books = append(books, RowToBook(b, db.Name))
	}

	return books, nil
}

func (db *DB) getAudiobookColumns() error {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	csql := sq.Case("label").
		When("'narrators'", "'books_custom_column_' || id || '_link'").
		Else("''")

	c, _, err := csql.ToSql()
	if err != nil {
		return err
	}

	key := sq.Select("id", "label", "name").
		Column(c + " 'join_table'").
		Column("'custom_column_' || id 'table'").
		From("custom_columns").
		Where(map[string]any{"label": []string{"narrators", "duration"}})

	stmt, args, err := key.ToSql()
	if err != nil {
		return err
	}

	rows, err := db.db.Queryx(stmt, args...)
	if err != nil {
		return fmt.Errorf("error %v\n, query\n %v", err, stmt)
	}
	defer rows.Close()
	db.db.Unsafe()

	for rows.Next() {
		var m Model
		err := rows.StructScan(&m)
		if err != nil {
			return err
		}
		mod := db.Models[m.Label]
		mod.JoinTable = m.JoinTable
		mod.Table = m.Table
		db.Models[m.Label] = mod
	}
	return nil
}

func ErrFileNotExist(path string) error {
	if !FileExist(path) {
		return fmt.Errorf("%v does not exist or cannot be found, check the path in the config, error: \n", path)
	}
	return nil
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
