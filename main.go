package main

import (
	"database/sql"
	"log"
	"runtime"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
	//"io/ioutil"
	"github.com/pavlik/fias_xml2postgresql/structures"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//var format = flag.String("format", "xml", "File format for import (xml or dbf)")
	//flag.Parse()

	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	structures.ExportActualStatus(dbmap)

}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", "user=dev dbname=fias password=dev sslmode=disable")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
