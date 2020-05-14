package db

import (
	"database/sql"
	"log"
	"sync"

	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	_ "github.com/mattn/go-sqlite3"
)

// lock is used to lock access to the sqlite file
var mutex = &sync.Mutex{}

type DB struct {
	Driver  string
	DataSrc string
	db      *sql.DB
}

func New(driver, datasource string) *DB {
	driver = dfaults.String("sqlite3", driver)
	return &DB{
		Driver:  driver,
		DataSrc: datasource,
	}
}

func (s *DB) Lock() {
	mutex.Lock()
}

func (s *DB) Unlock() {
	mutex.Unlock()
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

func (s *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := s.db.Exec(query, args...)
	if err != nil {
		log.Printf("Query error %s %s: %v", query, args, err)
		return nil, err
	}
	return result, nil
}
