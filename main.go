package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/structures/actual_status"
	"github.com/pavlik/fias_xml2postgresql/structures/address_object"
	"github.com/pavlik/fias_xml2postgresql/structures/center_status"
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

	var as_stat chan string = make(chan string, 1)
	var ao_stat chan string = make(chan string, 1)
	var cs_stat chan string = make(chan string, 1)

	if *format == "xml" {
		fmt.Println("обработка XML-файлов")
		go actual_status.Export(as_stat, db, format)
		go address_object.Export(ao_stat, db, format)
		go center_status.Export(cs_stat, db, format)
		// current_status.Export(dbmap)
		// estate_status.Export(dbmap)

	} else if *format == "dbf" {
		// todo: обработка DBF-файлов
		fmt.Println("обработка DBF-файлов")
	}

	var msg1, msg2, msg3 string
	for {
		select {
		case msg1 = <-as_stat:
		case msg2 = <-ao_stat:
		case msg3 = <-cs_stat:
		}
		progressPrint(msg1, msg2, msg3)
	}

	var input string
	fmt.Scanln(&input)
}

func progressPrint(msgs ...string) {
	var buffer bytes.Buffer

	for _, v := range msgs {
		buffer.WriteString(v)
		buffer.WriteString("\n")
	}

	fmt.Println(buffer.String())
	time.Sleep(time.Second * 1)
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
