package main

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
	"strings"

	"git.corvisa.com/uc/uc_migration/viewmodels"
)

func main() {

	// getting the pointer to templates files or cache
	templates := populateTemplates()

	http.HandleFunc("/",
		func(w http.ResponseWriter, req *http.Request) {
			requestedFile := req.URL.Path[1:]
			// matching requested url with templates. Need to call the templates as ex. ip:port/index
			template := templates.Lookup(requestedFile + ".html")

			var context interface{} = nil
			switch requestedFile {
			case "apps":
				context = viewmodels.GetApps(w, req)
			case "index":
				context = viewmodels.MigrateApps(w, req)
			}

			if template != nil {
				// here we inject data to the template as context object
				// we will get the data to store in viewmodel section of the project and then inject in the context object
				template.Execute(w, context)
			} else {
				w.WriteHeader(404)
			}
		})

	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)

	http.ListenAndServe(":9000", nil)
}

// serving static resource
func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "text/js"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else {
		contentType = "text/plain"
	}

	f, err := os.Open(path)

	if err == nil {
		defer f.Close()
		w.Header().Add("Content Type", contentType)

		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}

func populateTemplates() *template.Template {
	result := template.New("templates")

	basePath := "templates"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()

	templatePathsRaw, _ := templateFolder.Readdir(-1)

	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
				basePath+"/"+pathInfo.Name())
		}
	}

	result.ParseFiles(*templatePaths...)

	return result
}
