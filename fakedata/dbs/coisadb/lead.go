package coisadb

import (
	"database/sql"
	"fmt"

	"Fakedb/fakedata"

	"github.com/brianvoe/gofakeit/v6"
)

type FakeLeadTable struct {
	db *sql.DB
}

func NewFakeLeadTable(db *sql.DB) fakedata.FakeDataScript {
	return FakeLeadTable{db: db}
}

func (l FakeLeadTable) TableName() string {
	return "leads"
}

func (l FakeLeadTable) DB() *sql.DB {
	return l.db
}

func (l FakeLeadTable) FetchRows() (*sql.Rows, error) {
	return l.db.Query(fmt.Sprintf("SELECT id FROM %s", l.TableName()))
}

func (l FakeLeadTable) PrepareStmt(tx *sql.Tx) (*sql.Stmt, error) {
	updateQuery := fmt.Sprintf(`
		UPDATE %s
		SET email = ?
		WHERE id = ?;
	`, l.TableName())
	return tx.Prepare(updateQuery)
}

func (l FakeLeadTable) UpdateData(rows *sql.Rows, stmt *sql.Stmt) error {
	for rows.Next() {
		var ID int64
		err := rows.Scan(&ID)
		if err != nil {
			return fmt.Errorf("rows scan error: %w", err)
		}

		email := gofakeit.Email()

		_, err = stmt.Exec(email, ID)
		if err != nil {
			return fmt.Errorf("update exec error: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	return nil
}
