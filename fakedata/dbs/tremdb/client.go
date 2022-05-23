package tremdb

import (
	"database/sql"
	"fmt"

	"Fakedb/fakedata"
	"Fakedb/fakedata/helper"
)

type FakeClientTable struct {
	db *sql.DB
}

func NewFakeClientTable(db *sql.DB) fakedata.FakeDataScript {
	return FakeClientTable{db: db}
}

func (c FakeClientTable) TableName() string {
	return "client"
}

func (c FakeClientTable) DB() *sql.DB {
	return c.db
}

func (c FakeClientTable) FetchRows() (*sql.Rows, error) {
	return c.db.Query(fmt.Sprintf("SELECT id FROM %s", c.TableName()))
}

func (c FakeClientTable) PrepareStmt(tx *sql.Tx) (*sql.Stmt, error) {
	updateQuery := fmt.Sprintf(`
		UPDATE %s
		SET name = ?, cpf = ?, address = ?
		WHERE id = ?;
	`, c.TableName())
	return tx.Prepare(updateQuery)
}

func (c FakeClientTable) UpdateData(rows *sql.Rows, stmt *sql.Stmt) error {
	for rows.Next() {
		var ID int64
		err := rows.Scan(&ID)
		if err != nil {
			return fmt.Errorf("rows scan error: %w", err)
		}

		name := helper.GenerateName()
		cpf := helper.GenerateCpf()
		address := helper.GenerateAddress()

		_, err = stmt.Exec(name, cpf, address, ID)
		if err != nil {
			return fmt.Errorf("update exec error: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	return nil
}
