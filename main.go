package main

import (
	"io/ioutil"
	"net/http"
)

func main() {

	// http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	w.Write([]byte("Hello World"))
	// })

	http.Handle("/", new(MyHandler))

	http.ListenAndServe(":8000", nil)
}

type MyHandler struct {
	http.Handler // we will implement this handler
}

func (this *MyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path

	// ReadFile will load all the resources in the path to the Memory, so better to use bufferedReader
	data, err := ioutil.ReadFile(string(path))

	if err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
	}

}
