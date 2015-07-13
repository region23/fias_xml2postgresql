package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/structures/actual_status"
	"github.com/pavlik/fias_xml2postgresql/structures/address_object"
	"github.com/pavlik/fias_xml2postgresql/structures/center_status"
	"github.com/pavlik/fias_xml2postgresql/structures/current_status"
	"github.com/pavlik/fias_xml2postgresql/structures/estate_status"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("Используемое количество ядер: %d\n", runtime.NumCPU())

	var format = flag.String("format", "xml", "File format for import (xml or dbf)")
	flag.Parse()

	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	//progress := make(chan string)

	if *format == "xml" {
		fmt.Println("обработка XML-файлов")
		actual_status.Export(dbmap)
		address_object.Export(dbmap)
		center_status.Export(dbmap)
		current_status.Export(dbmap)
		estate_status.Export(dbmap)

	} else if *format == "dbf" {
		// todo: обработка DBF-файлов
		fmt.Println("обработка DBF-файлов")
	}

	// var input string
	// fmt.Scanln(&input)
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
