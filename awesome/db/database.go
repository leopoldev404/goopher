package db

import (
	"awesome/types"
	"database/sql"
	"fmt"
)

type mysqlDb struct {
	db *sql.DB
}

func NewMySqlDb(connectionString string) (*mysqlDb, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	return &mysqlDb{db: db}, nil
}

func (db *mysqlDb) Save(user *types.User) error {
	return nil
}

func (db *mysqlDb) GetById(id string) (*types.User, error) {
	return nil, nil
}

func (db *mysqlDb) Get(username, password string) (*types.User, error) {
	return nil, nil
}

func (db *mysqlDb) Dispose() {
	db.db.Close()
}
