package script

import (
	"Fakedb/credentials"
	"Fakedb/fakedata"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func FakeData(fd fakedata.IFakeData, cp credentials.ConnectionParams) error {

	db, err := getDBHandle(cp)
	if err != nil {
		return err
	}
	defer db.Close()

	err = runFakeDataScripts(fd.Scripts(db))
	if err != nil {
		return err
	}

	return nil
}

func getDBHandle(credentials credentials.ConnectionParams) (*sql.DB, error) {
	addr := fmt.Sprintf("%s:%s", credentials.Host, credentials.Port)
	cfg := mysql.Config{
		User:                 credentials.User,
		Passwd:               credentials.Password,
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               credentials.Database,
		AllowNativePasswords: true,
	}

	return getConnection(cfg)
}

func getConnection(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		err = fmt.Errorf("erro to open database: %w", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		err = fmt.Errorf("erro to ping database: %w", err)
		return nil, err
	}

	return db, nil
}

func runFakeDataScripts(scripts []fakedata.FakeDataScript) error {

	for _, s := range scripts {
		rows, err := s.FetchRows()
		if err != nil {
			return err
		}
		defer rows.Close()

		tx, err := s.DB().Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		stmt, err := s.PrepareStmt(tx)
		if err != nil {
			return err
		}
		defer stmt.Close()

		err = s.UpdateData(rows, stmt)
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}
