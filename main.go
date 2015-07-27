package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nsf/termbox-go"

	_ "github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/structures/actual_status"
	"github.com/pavlik/fias_xml2postgresql/structures/address_object"
	"github.com/pavlik/fias_xml2postgresql/structures/center_status"
	"github.com/pavlik/fias_xml2postgresql/structures/current_status"
	"github.com/pavlik/fias_xml2postgresql/structures/estate_status"
	"github.com/pavlik/fias_xml2postgresql/structures/house"
)

func print_tb(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func printf_tb(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	print_tb(x, y, fg, bg, s)
}

func progressPrint(msgs ...string) {
	color := [8]termbox.Attribute{termbox.ColorRed, termbox.ColorGreen, termbox.ColorYellow, termbox.ColorBlue, termbox.ColorMagenta, termbox.ColorCyan, termbox.ColorWhite}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	printf_tb(0, 0, termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack, "Экспорт базы ФИАС в БД PostgreSQL")

	y := 0
	for _, v := range msgs {
		printf_tb(0, y+2, color[y], termbox.ColorBlack, v)
		y++
	}

	termbox.Flush()
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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf("Используемое количество ядер: %d\n", runtime.NumCPU())

	var format = flag.String("format", "xml", "File format for import (xml or dbf)")
	flag.Parse()

	// initialize the DbMap
	db := initDb()
	defer db.Close()

	var as_stat chan string = make(chan string, 1000)
	var ao_stat chan string = make(chan string, 1000)
	var cs_stat chan string = make(chan string, 1000)
	var cur_stat chan string = make(chan string, 1000)
	var est_stat chan string = make(chan string, 1000)
	var house_stat chan string = make(chan string, 1000)

	if *format == "xml" {
		fmt.Println("обработка XML-файлов")
		go actual_status.Export(as_stat, db, format)
		go address_object.Export(ao_stat, db, format)
		go center_status.Export(cs_stat, db, format)
		go current_status.Export(cur_stat, db, format)
		go estate_status.Export(est_stat, db, format)
		go house.Export(house_stat, db, format)

	} else if *format == "dbf" {
		// todo: обработка DBF-файлов
		fmt.Println("обработка DBF-файлов")
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	var msg1, msg2, msg3, msg4, msg5, msg6 string
	timer := time.After(time.Second * 1)
	go func() {
		for {
			select {
			case <-timer:
				progressPrint(msg1, msg2, msg3, msg4, msg5, msg6)
				timer = time.After(time.Second * 1)
			default:
			}
			select {
			case msg1 = <-as_stat:
			case msg2 = <-ao_stat:
			case msg3 = <-cs_stat:
			case msg4 = <-cur_stat:
			case msg5 = <-est_stat:
			case msg6 = <-house_stat:
			}
		}
	}()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				break loop
			}
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			termbox.Flush()
			//progressPrint(msg1, msg2, msg3, msg4, msg5)
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
