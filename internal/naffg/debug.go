package naffg

import (
	"fmt"
	"path/filepath"
)

func (app *Application) Debug_PrintEmbeddedFiles() {
	if app.options.UiRes != nil {
		fmt.Println("Loaded Embedded Files:")
		var printFiles func(dir string)
		printFiles = func(dir string) {
			fileList, _ := app.options.UiRes.ReadDir(dir)
			for _, file := range fileList {
				fullPath := filepath.Join(dir, file.Name())
				fmt.Println(fullPath)
				if file.IsDir() {
					printFiles(fullPath)
				}
			}
		}
		printFiles(".")
	}
}
