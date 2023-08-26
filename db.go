package cdb

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	Name         string
	Models       Models
	db           *sqlx.DB
	mtx          sync.Mutex
	printQuery   bool
	isAudiobooks bool
	sort         string
	limit        int
	query        sq.SelectBuilder
}

func configDB(name, path string) (*DB, error) {
	db := &DB{
		Name:   name,
		Models: modelMeta,
		sort:   "timestamp",
	}

	p := filepath.Join(path, name, metaDB)
	if ok := FileExist(p); !ok {
		return db, ErrFileNotExist(p)
	}

	url := sqlitePrefix + p + sqliteOpts
	database, err := sqlx.Open("sqlite3", url)
	if err != nil {
		return db, fmt.Errorf("database connection %v failed\n", err)
	}

	db.db = database

	return db, nil
}

func (db DB) IsConnected() bool {
	return db.db != nil
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
		b := &Book{}
		err := rows.StructScan(b)
		if err != nil {
			return books, err
		}
		books = append(books, b)
	}

	return books, nil
}

func (db *DB) getAudiobookColumns() error {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	db.Models = AudiobookModels()

	csql := sq.Case("label").
		When("'narrators'", "'books_custom_column_' || id || '_link'").
		Else("''")

	c, _, err := csql.ToSql()
	if err != nil {
		return err
	}

	key := sq.Select("id", "'#' || label 'label'", "name").
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
