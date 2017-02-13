package main

import "net/http"

func main() {

	// built in server to handle request
	http.ListenAndServe(":9000", http.FileServer(http.Dir("public")))
}

type MyHandler struct {
	http.Handler // we will implement this handler
}
