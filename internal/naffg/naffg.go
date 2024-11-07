package naffg

/*
	Native-like Application Framework For Go (NAFFG)
*/

import (
	"embed"
	"net/http"
	"strconv"

	webview2 "github.com/jchv/go-webview2"
)

type Application struct {
	viewport     webview2.WebView
	options      *Options
	backendMutex *http.ServeMux
	terminate    chan bool
}

type Options struct {
	//Basic Settings
	Title            string
	Width            int
	Height           int
	Resizable        bool //Default is true
	InitiateAtCenter bool //Initiate the window at the center of the screen

	//WebApp Settings
	StartURI string //URI to be loaded in the webview on start, support data://, default is "http://localhost:36850"

	//Local App Settings
	UiRes      *embed.FS //Go embed filesystem for UI resources (e.g. res/), leave empty to use external URI (webApp)
	UiResPreix string    //Prefix for the UI resources, will trim the prefix before serving the resources (e.g. res/)

	//Advance Settings
	EventExchangePort int //Port for which the UI events will be exchanged with backend, default is 36850
	IconId            int //Icon resource id
	Debug             bool
}

// Create a new Application Window with the given options
func NewApplication(options *Options) *Application {
	//Fill in the default values if not provided
	if options.Title == "" {
		options.Title = "New Application"
	}

	if options.Width == 0 {
		options.Width = 800
	}

	if options.Height == 0 {
		options.Height = 600
	}

	if options.EventExchangePort == 0 {
		options.EventExchangePort = 36850
	}

	if options.UiRes == nil && options.StartURI == "" {
		//Use a missing resource page
		options.StartURI = "data:text/html,<html><body><h1>Invalid Usage</h1></body></html>"
	}

	//Create the webview
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  options.Title,
			Width:  uint(options.Width),
			Height: uint(options.Height),
			IconId: 203, // icon resource id
			Center: options.InitiateAtCenter,
		},
	})

	//Create a http mux for backend
	mux := http.NewServeMux()

	app := &Application{
		viewport:     w,
		options:      options,
		backendMutex: mux,
		terminate:    make(chan bool),
	}

	//Register the UI resources bundle if provided
	if options.UiRes != nil {
		mux.Handle("/", app.embedFsPrefixMiddleware(http.FileServer(http.FS(options.UiRes))))
	}
	return app
}

// Run the application
// Note that this cannot be run in a goroutine. It will block the main thread
// until the application is terminated
func (app *Application) Run() {
	app.viewport.SetTitle(app.options.Title)
	if app.options.Resizable {
		app.viewport.SetSize(app.options.Width, app.options.Height, webview2.HintNone)
	} else {
		app.viewport.SetSize(app.options.Width, app.options.Height, webview2.HintFixed)
	}

	if app.options.StartURI != "" {
		app.viewport.Navigate(app.options.StartURI)
	} else {
		app.viewport.Navigate("http://localhost:36850")
	}

	//Start the backend server
	go func() {
		if err := http.ListenAndServe("localhost:"+strconv.Itoa(app.options.EventExchangePort), app.backendMutex); err != nil {
			panic(err)
		}
	}()

	//Run the webview
	defer app.viewport.Destroy()
	app.viewport.Run()

}

// Terminate the application
func (app *Application) Terminate() {
	app.viewport.Destroy()
}
