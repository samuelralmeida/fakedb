package fakedata

import "database/sql"

type FakeDataScript interface {
	DB() *sql.DB
	FetchRows() (*sql.Rows, error)
	PrepareStmt(tx *sql.Tx) (*sql.Stmt, error)
	UpdateData(*sql.Rows, *sql.Stmt) error
	TableName() string
}

type IFakeData interface {
	Scripts(db *sql.DB) []FakeDataScript
}
