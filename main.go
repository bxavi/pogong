package main

import (
	"database/sql"

	"github.com/bxavi/pogong/api"
	"github.com/bxavi/pogong/db"
	_ "github.com/lib/pq"

	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:password@localhost:5432/pogong?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	con, err := sql.Open(dbDriver, dbSource)
	if err != nil || con == nil {
		log.Fatal("TestMain no access to be db")
	}

	store := db.NewStore(con)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
