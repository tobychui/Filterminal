package naffg

import (
	"net/http"
	"path/filepath"
)

// Mux returns the backend mutex for user to register their own handlers
func (app *Application) Mux() *http.ServeMux {
	return app.backendMutex
}

func (app *Application) embedFsPrefixMiddleware(fsHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/index.html", http.StatusFound)
			return
		}

		if app.options.UiResPreix != "" {
			r.URL.Path = filepath.Join(app.options.UiResPreix, r.URL.Path)
			//Force override of mime type if html, js or css
			if filepath.Ext(r.URL.Path) == ".html" {
				w.Header().Set("Content-Type", "text/html")
			} else if filepath.Ext(r.URL.Path) == ".js" {
				w.Header().Set("Content-Type", "application/javascript")
			} else if filepath.Ext(r.URL.Path) == ".css" {
				w.Header().Set("Content-Type", "text/css")
			}

			if app.options.Debug {
				println("Serving: ", r.URL.Path)
			}
			http.ServeFile(w, r, r.URL.Path)
			return
		}

		fsHandler.ServeHTTP(w, r)
	})

}
