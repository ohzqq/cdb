package cdb

import (
	"errors"
	"fmt"
	"os"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqliteOpts       = "?cache=shared&mode=ro"
	sqlitePrefix     = `file:`
	metaDB           = `metadata.db`
	SavedSearches    = "saved_searches"
	DisplayFields    = "book_display_fields"
	HiddenCategories = "tag_browser_hidden_categories"
	FieldMetadata    = "field_metadata"
)

// DB holds the database connection.
type DB struct {
	db         *sqlx.DB
	printQuery bool
	mtx        sync.Mutex
}

// Connect opens the database connection.
func (db *DB) Connect(p string) error {
	if db.IsConnected() {
		return nil
	}

	if ok := FileExist(p); !ok {
		return ErrFileNotExist(p)
	}

	url := sqlitePrefix + p + sqliteOpts
	database, err := sqlx.Open("sqlite3", url)
	if err != nil {
		return fmt.Errorf("database connection %v failed\n", err)
	}

	db.db = database
	return nil
}

// IsConnected checks if the database is connected.
func (db DB) IsConnected() bool {
	return db.db != nil
}

func (db *DB) bookQuery(stmt string, args []any) (Records, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if db.printQuery {
		fmt.Println(stmt)
	}
	var records Records
	rows, err := db.db.Queryx(stmt, args...)
	if err != nil {
		fmt.Println(stmt)
		return records, fmt.Errorf("error %v\n", err)
	}
	defer rows.Close()
	db.db.Unsafe()

	for rows.Next() {
		meta := make(map[string]any)
		scanErr := rows.MapScan(meta)
		if scanErr != nil {
			return records, scanErr
		}
		records.rows = append(records.rows, meta)
	}
	return records, nil
}

func (db *DB) getAudiobookColumns() (Models, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	models := AudiobookModels()

	csql := sq.Case("label").
		When("'narrators'", "'books_custom_column_' || id || '_link'").
		Else("''")

	c, _, err := csql.ToSql()
	if err != nil {
		return models, err
	}

	key := sq.Select("id", "label", "name").
		Column(c + " 'join_table'").
		Column("'custom_column_' || id 'table'").
		From("custom_columns").
		Where(map[string]any{"label": []string{"narrators", "duration"}})

	stmt, args, err := key.ToSql()
	if err != nil {
		return models, err
	}

	rows, err := db.db.Queryx(stmt, args...)
	if err != nil {
		return models, fmt.Errorf("error %v\n, query\n %v", err, stmt)
	}
	defer rows.Close()
	db.db.Unsafe()

	for rows.Next() {
		var m Model
		err := rows.StructScan(&m)
		if err != nil {
			return models, err
		}
		mod := models[m.Label]
		mod.JoinTable = m.JoinTable
		mod.Table = m.Table
		models[m.Label] = mod
	}
	return models, nil
}

func (db *DB) getPreferences(stmt string, args ...any) ([]byte, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	var data []byte
	row := db.db.QueryRowx(stmt, args...)
	row.Scan(&data)

	return data, nil
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
