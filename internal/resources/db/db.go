package db

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// mu Mutex used to co-ordinate all writes to a sequential process
// instead of parallel as per Sqlite3 documentation.
var mu = &sync.Mutex{}

type DB struct {
	db *sql.DB
}

// Options provides Options to New
type Options struct {
	DB *sql.DB
}

// New Creates a *DB instance. Driver defaults to "sqlite3" if driver is an empty string.
func New(opts *Options) *DB {
	if opts.DB == nil {
		panic("error Options.DB can not be nil")
	}
	return &DB{
		db: opts.DB,
	}
}

func (s *DB) Lock() {
	Lock()
}

func Lock() {
	mu.Lock()
}

func (s *DB) Unlock() {
	Unlock()
}

func Unlock() {
	mu.Unlock()
}

func (s *DB) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	result, err = s.db.Exec(query, args...)

	if err != nil {
		log.Printf("Query error %s %s: %v", query, args, err)
		return nil, err
	}
	return result, nil
}
