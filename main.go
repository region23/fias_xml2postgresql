package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nsf/termbox-go"

	_ "github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/structures/actual_status"
	"github.com/pavlik/fias_xml2postgresql/structures/address_object"
	"github.com/pavlik/fias_xml2postgresql/structures/address_object_type"
	"github.com/pavlik/fias_xml2postgresql/structures/center_status"
	"github.com/pavlik/fias_xml2postgresql/structures/current_status"
	"github.com/pavlik/fias_xml2postgresql/structures/estate_status"
	"github.com/pavlik/fias_xml2postgresql/structures/house"
	"github.com/pavlik/fias_xml2postgresql/structures/house_interval"
	"github.com/pavlik/fias_xml2postgresql/structures/house_state_status"
	"github.com/pavlik/fias_xml2postgresql/structures/interval_status"
	"github.com/pavlik/fias_xml2postgresql/structures/landmark"
	"github.com/pavlik/fias_xml2postgresql/structures/normative_document"
	"github.com/pavlik/fias_xml2postgresql/structures/normative_document_type"
	"github.com/pavlik/fias_xml2postgresql/structures/operation_status"
	"github.com/pavlik/fias_xml2postgresql/structures/structure_status"
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

const timeLayout = "2006-01-02 в 15:04"

func progressPrint(msgs [15]string, counters [15]int, startTime time.Time, finished bool) {
	color := [15]termbox.Attribute{termbox.ColorWhite,
		termbox.ColorWhite,
		termbox.ColorWhite,
		termbox.ColorWhite,
		termbox.ColorWhite,
		termbox.ColorBlue,
		termbox.ColorBlue,
		termbox.ColorBlue,
		termbox.ColorBlue,
		termbox.ColorBlue,
		termbox.ColorRed,
		termbox.ColorRed,
		termbox.ColorRed,
		termbox.ColorRed,
		termbox.ColorRed}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	printf_tb(0, 0, termbox.ColorGreen|termbox.AttrBold, termbox.ColorBlack, "Экспорт базы ФИАС в БД PostgreSQL")
	printf_tb(0, 1, termbox.ColorYellow, termbox.ColorBlack, fmt.Sprintf("Количество используемых ядер: %d", runtime.NumCPU()))
	printf_tb(0, 20, termbox.ColorCyan, termbox.ColorBlack, fmt.Sprintf("Конвертация началась %s", startTime.Format(timeLayout)))

	duration := time.Since(startTime)
	var durationText string

	if duration.Seconds() < 60 {
		durationText = fmt.Sprintf("%.1f секунд", duration.Seconds())
	} else if duration.Minutes() < 60 {
		durationText = fmt.Sprintf("%.1f минут", duration.Minutes())
	} else {
		durationText = fmt.Sprintf("%.1f часов", duration.Hours())
	}

	if finished {
		printf_tb(0, 21, termbox.ColorGreen, termbox.ColorBlack, fmt.Sprintf("База экспортировалась %s. Экспорт завершен.", durationText))
	} else {
		printf_tb(0, 21, termbox.ColorCyan, termbox.ColorBlack, fmt.Sprintf("и уже длится %s", durationText))
	}
	printf_tb(0, 22, termbox.ColorMagenta|termbox.AttrUnderline, termbox.ColorBlack, "Для прерывания экспорта и выхода из программы нажмите CTRL+Q")

	y := 0
	for _, v := range msgs {
		if counters[y] > 0 {
			v = fmt.Sprintf("%s. Total count is %d", v, counters[y])
		}
		printf_tb(0, y+3, color[y], termbox.ColorDefault, v)
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

func logInit() *log.Logger {
	logFile := "log.txt"

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", logFile, ":", err)
	}

	logger := log.New(file,
		"FATAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return logger
}

func main() {

	logger := logInit()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Program terminated. Main goroutine paniced:", r)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	var format = flag.String("format", "xml", "File format for import (xml or dbf)")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// initialize the DbMap
	db := initDb()
	defer db.Close()

	var w sync.WaitGroup
	w.Add(15)

	as_stat := make(chan string, 1000)
	ao_stat := make(chan string, 1000)
	//ao_counter := make(chan int, 1000)
	cs_stat := make(chan string, 1000)
	cur_stat := make(chan string, 1000)
	est_stat := make(chan string, 1000)
	house_stat := make(chan string, 1000)
	//house_counter := make(chan int, 1000)
	house_int_stat := make(chan string, 1000)
	//house_int_counter := make(chan int, 1000)
	house_st_stat := make(chan string, 1000)
	intv_stat := make(chan string, 1000)
	landmark_stat := make(chan string, 1000)
	//landmark_counter := make(chan int, 1000)
	ndtype_stat := make(chan string, 1000)
	nd_stat := make(chan string, 1000)
	//nd_counter := make(chan int, 1000)
	oper_stat := make(chan string, 1000)
	socrbase_stat := make(chan string, 1000)
	str_stat := make(chan string, 1000)

	done := make(chan bool, 1)
	//done <- false
	//logger.Println(<-done)

	if *format == "xml" {
		fmt.Println("обработка XML-файлов")

		go actual_status.ExportBulk(&w, as_stat, db, format, logger)
		go estate_status.ExportBulk(&w, est_stat, db, format, logger)
		go interval_status.ExportBulk(&w, intv_stat, db, format, logger)
		go structure_status.ExportBulk(&w, str_stat, db, format, logger)
		go center_status.ExportBulk(&w, cs_stat, db, format, logger)

		go operation_status.ExportBulk(&w, oper_stat, db, format, logger)
		go normative_document_type.ExportBulk(&w, ndtype_stat, db, format, logger)
		go house_state_status.ExportBulk(&w, house_st_stat, db, format, logger)
		go current_status.ExportBulk(&w, cur_stat, db, format, logger)
		go address_object_type.ExportBulk(&w, socrbase_stat, db, format, logger)

		go landmark.ExportBulk(&w, landmark_stat, db, format, logger)
		//go helpers.CountElementsInXML(&w, landmark_counter, "as_landmark", "Landmark")

		go normative_document.ExportBulk(&w, nd_stat, db, format, logger)
		//go helpers.CountElementsInXML(&w, nd_counter, "as_normdoc", "NormativeDocument")

		go house_interval.ExportBulk(&w, house_int_stat, db, format, logger)
		//go helpers.CountElementsInXML(&w, house_int_counter, "as_houseint", "HouseInterval")

		go address_object.ExportBulk(&w, ao_stat, db, format, logger)
		//go helpers.CountElementsInXML(&w, ao_counter, "as_addrobj", "Object")

		go house.ExportBulk(&w, house_stat, db, format, logger)
		//go helpers.CountElementsInXML(&w, house_counter, "as_house_", "House")

	} else if *format == "dbf" {
		// todo: обработка DBF-файлов
		fmt.Println("обработка DBF-файлов")
	}

	err := termbox.Init()
	if err != nil {
		logger.Fatal(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	var msgs [15]string
	var counters [15]int

	doneStat := false

	timer := time.After(time.Second * 1)
	go func() {

		startTime := time.Now()

		//	progressLoop:
		for {
			select {
			case doneStat = <-done:
				logger.Println("Попали в done")
				progressPrint(msgs, counters, startTime, true)

				//break progressLoop
				return
			case <-timer:
				if doneStat {
					logger.Println("Попали в таймер после done")
				}
				progressPrint(msgs, counters, startTime, false)
				timer = time.After(time.Millisecond)
			// default:
			// }
			// select {
			case msgs[0] = <-as_stat:
				if doneStat {
					logger.Println("Попали в msgs[0] после done")
				}
			case msgs[1] = <-est_stat:
				if doneStat {
					logger.Println("Попали в msgs[1] после done")
				}
			case msgs[2] = <-intv_stat:
				if doneStat {
					logger.Println("Попали в msgs[2] после done")
				}
			case msgs[3] = <-str_stat:
				if doneStat {
					logger.Println("Попали в msgs[3] после done")
				}
			case msgs[4] = <-cs_stat:
				if doneStat {
					logger.Println("Попали в msgs[4] после done")
				}
			case msgs[5] = <-oper_stat:
				if doneStat {
					logger.Println("Попали в msgs[5] после done")
				}
			case msgs[6] = <-ndtype_stat:
				if doneStat {
					logger.Println("Попали в msgs[6] после done")
				}
			case msgs[7] = <-house_st_stat:
				if doneStat {
					logger.Println("Попали в msgs[7] после done")
				}
			case msgs[8] = <-cur_stat:
				if doneStat {
					logger.Println("Попали в msgs[8] после done")
				}
			case msgs[9] = <-socrbase_stat:
				if doneStat {
					logger.Println("Попали в msgs[9] после done")
				}
			case msgs[10] = <-landmark_stat:
				if doneStat {
					logger.Println("Попали в msgs[10] после done")
				}
			case msgs[11] = <-nd_stat:
				if doneStat {
					logger.Println("Попали в msgs[11] после done")
				}
			case msgs[12] = <-house_int_stat:
				if doneStat {
					logger.Println("Попали в msgs[12] после done")
				}
			case msgs[13] = <-ao_stat:
				if doneStat {
					logger.Println("Попали в msgs[13] после done")
				}
			case msgs[14] = <-house_stat:
				if doneStat {
					logger.Println("Попали в msgs[14] после done")
				}
			}
			// select {
			// // case msgs[0] = <-as_stat:
			// // case msgs[1] = <-est_stat:
			// // case msgs[2] = <-intv_stat:
			// // case msgs[3] = <-str_stat:
			// // case msgs[4] = <-cs_stat:
			// // case msgs[5] = <-oper_stat:
			// // case msgs[6] = <-ndtype_stat:
			// // case msgs[7] = <-house_st_stat:
			// // case msgs[8] = <-cur_stat:
			// // case msgs[9] = <-socrbase_stat:
			// // case counters[10] = <-landmark_counter:
			// // case counters[11] = <-nd_counter:
			// // case counters[12] = <-house_int_counter:
			// // case counters[13] = <-ao_counter:
			// // case counters[14] = <-house_counter:
			// }
		}
	}()

	go func() {
		w.Wait()
		done <- true
		logger.Println("Все горутины завершились")
		// или close(done), по желанию
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
			logger.Fatal(err)
		}
	}

}
