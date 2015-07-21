package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/structures/actual_status"
	"github.com/pavlik/fias_xml2postgresql/structures/address_object"
	// "github.com/pavlik/fias_xml2postgresql/structures/center_status"
	// "github.com/pavlik/fias_xml2postgresql/structures/current_status"
	// "github.com/pavlik/fias_xml2postgresql/structures/estate_status"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("Используемое количество ядер: %d\n", runtime.NumCPU())

	var format = flag.String("format", "xml", "File format for import (xml or dbf)")
	flag.Parse()

	// initialize the DbMap
	db := initDb()
	defer db.Close()

	if *format == "xml" {
		fmt.Println("обработка XML-файлов")
		actual_status.Export(db, format)
		address_object.Export(db, format)
		// center_status.Export(dbmap)
		// current_status.Export(dbmap)
		// estate_status.Export(dbmap)

	} else if *format == "dbf" {
		// todo: обработка DBF-файлов
		fmt.Println("обработка DBF-файлов")
	}

	// var input string
	// fmt.Scanln(&input)
}

func initDb() *sqlx.DB {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sqlx.Open("postgres", "user=dev dbname=fias password=dev sslmode=disable")
	checkErr(err, "sqlx.Open failed")

	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
