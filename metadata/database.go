package metadata

import (
	"github.com/prologic/bitcask"
	"kproxy/helpers"
)

var _db *bitcask.Bitcask

func GetDatabaseSingleton() *bitcask.Bitcask {
	if _db == nil {
		panic("Database not yet initialised")
	}

	return _db
}

func Init() {
	db, err := bitcask.Open(helpers.GetDatabasePath())
	if err != nil {
		panic(err)
	}

	_db = db
}
