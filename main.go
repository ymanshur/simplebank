package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ymanshur/simplebank/util"

	"github.com/ymanshur/simplebank/api"
	db "github.com/ymanshur/simplebank/db/sqlc"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)

	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
