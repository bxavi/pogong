package main

import (
	"database/sql"

	"github.com/bxavi/pogong/api"
	"github.com/bxavi/pogong/db"
	"github.com/bxavi/pogong/util"
	_ "github.com/lib/pq"

	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config file:", err)
	}

	con, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil || con == nil {
		log.Fatal("no access to be db:", err)
	}

	store := db.NewStore(con)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create new server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
