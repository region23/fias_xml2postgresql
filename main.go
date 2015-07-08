package main

import (
	"flag"
	"runtime"
	//"io/ioutil"
	"github.com/pavlik/fias_xml2postgres/structures"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var format = flag.String("format", "xml", "File format for import (xml or dbf)")
	flag.Parse()

	structures.ExportActualStatus()

}
