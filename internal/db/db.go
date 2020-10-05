package db

import (
	"database/sql"
	"errors"
	"log"
	"sync"

	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	_ "github.com/mattn/go-sqlite3"
)

// mu Mutex used to co-ordinate all writes to a sequential process
// instead of parallel as per Sqlite3 documentation.
var mu = &sync.Mutex{}

type DB struct {
	Driver  string
	DataSrc string
	db      *sql.DB
}

// New Creates a *DB instance. Driver defaults to "sqlite3" if driver is an empty string.
func New(driver, datasource string) *DB {
	driver = dfaults.String("sqlite3", driver)
	if datasource == "" {
		errutils.Epanicf("ERROR: %w", errors.New("Datasource is empty"))
	}

	return &DB{
		Driver:  driver,
		DataSrc: datasource,
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

func (s *DB) Open() *sql.DB {
	if s.db != nil {
		return s.db
	}
	db, err := sql.Open("sqlite3", s.DataSrc)
	s.db = db
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}

func (s *DB) Close() {
	if s.db != nil {
		s.db.Close()
		s.db = nil
	}
}

func (s *DB) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	result, err = s.db.Exec(query, args...)

	if err != nil {
		log.Printf("Query error %s %s: %v", query, args, err)
		return nil, err
	}
	return result, nil
}
