package main

import (
	"github.com/gin-gonic/gin"
	"entrega/internal/database"
	"entrega/internal/service/avocados"
	"flag"
	"os"
	"fmt"
	"entrega/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

)

func main() {
	cfg:= readConfig()

	db, err := database.NewDatabase(cfg)
	defer db.Close()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = createSchema(db)
	if err!= nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	service, _:= avocados.New(db, cfg)
	httpService:= avocados.NewHTTPTransport(service)
	r:= gin.Default()
	httpService.Register(r)
	r.Run()
}

func readConfig() *config.Config {
	configFile := flag.String("config", "./config.yaml", "service config")
	flag.Parse()
	cfg, err := config.LoadConfig(*configFile)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return cfg
}

func createSchema(db *sqlx.DB) error {
	schema :=  `CREATE TABLE IF NOT EXISTS avocados (
				id integer primary key autoincrement,
				name varchar,
				image varchar,
				price integer,
				stock integer,
				quantity integer);`
		_, err:= db.Exec(schema)

	if err != nil {
		return err
	}
	return nil
}
