package viewmodels

import (
	"fmt"
	"net/http"
)

type Apps struct {
	Source string
	Target string
}

func GetApps(w http.ResponseWriter, req *http.Request) Apps {

	err := req.ParseForm()
	if err != nil {
		// Handle error here via logging and then return
	}

	src := req.PostFormValue("source")
	tar := req.PostFormValue("target")
	fmt.Printf("src : %s", src)
	result := Apps{
		Source: src,
		Target: tar,
	}

	return result
}
