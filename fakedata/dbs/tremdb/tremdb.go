package tremdb

import (
	"Fakedb/fakedata"
	"database/sql"
)

type TremDB struct{}

func NewTremDb() fakedata.IFakeData {
	return TremDB{}
}

func (c TremDB) Scripts(db *sql.DB) []fakedata.FakeDataScript {
	return []fakedata.FakeDataScript{
		NewFakeClientTable(db),
	}
}
