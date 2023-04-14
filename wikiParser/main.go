package main

import (
	"log"
	"os"
	"wiki-parser/data"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	l := log.New(os.Stdout, "wikiParser: ", log.LstdFlags)
	dbProvider, err := data.NewDbProvider(l)
	if err != nil {
		l.Fatal("Cannot connect to database")
	}

	dbProvider.ParsePages()
	defer dbProvider.Close()
}
