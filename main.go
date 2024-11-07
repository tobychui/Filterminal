package main

import (
	"embed"
	"flag"
	"net/http"
	"strconv"

	"imuslab.com/aroz/filterm/internal/naffg"
)

/*
	Filterminal
*/

var (
	//Global variables
	applicationInstance *naffg.Application

	//System information
	version = "0.0.1"

	//Startup flags
	comport = flag.String("comport", "", "COM port to use for communication")
	baud    = flag.Int("baud", 115200, "Baud rate for communication")
	help    = flag.Bool("help", false, "Show this help message")
)

// Embed ui resources

//go:embed res/*
var uiResources embed.FS

func main() {
	flag.Parse()

	//Start naffg interface
	applicationInstance = naffg.NewApplication(&naffg.Options{
		Title:            "Filterminal " + version,
		Width:            800,
		Height:           600,
		Resizable:        true,
		InitiateAtCenter: true,
		UiRes:            &uiResources,
		UiResPreix:       "res",
	})

	applicationInstance.Mux().HandleFunc("/comport", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(*comport))
	})

	applicationInstance.Mux().HandleFunc("/baud", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.Itoa(*baud)))
	})

	applicationInstance.Run()
}
