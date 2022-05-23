package coisadb

import (
	"Fakedb/fakedata"
	"database/sql"
)

type CoisaDb struct{}

func NewCoisaDb() fakedata.IFakeData {
	return CoisaDb{}
}

func (c CoisaDb) Scripts(db *sql.DB) []fakedata.FakeDataScript {
	return []fakedata.FakeDataScript{
		NewFakeLeadTable(db),
	}
}
