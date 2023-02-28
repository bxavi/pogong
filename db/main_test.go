package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/bxavi/pogong/util"
	_ "github.com/lib/pq"
)

var testDb *sql.DB
var testQueries *Queries = New(testDb)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load confid:", err)
	}
	con, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil || con == nil {
		log.Fatal("TestMain no access to be db")
	}
	testQueries = New(con)
	os.Exit(m.Run())
}
