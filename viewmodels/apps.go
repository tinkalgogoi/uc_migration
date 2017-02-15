package viewmodels

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var Apps map[string]interface{}

func GetApps(w http.ResponseWriter, req *http.Request) []map[string]interface{} {
	marathonIP := "http://10.198.161.41:8080/v2/apps"

	m := make(map[string]interface{})

	err := req.ParseForm()
	if err != nil {
		fmt.Printf("ParseForm err : %s", err)
	}

	src := req.PostFormValue("source")
	//tar := req.PostFormValue("target")
	// put some strictly matching query for the rest call
	marathonURL := marathonIP + "?id=" + src
	res, err := http.Get(marathonURL)
	if err != nil {
		fmt.Printf("marathonURL err : %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("err : %s", err)
	}

	json.Unmarshal(body, &Apps)
	//fmt.Println(Apps)
	var appsSlice []map[string]interface{}
	// Getting each map at the apps env field
	for i := range Apps["apps"].([]interface{}) {
		// fmt.Println(get(Apps, "apps", i, "id"))
		// fmt.Println(get(Apps, "apps", i, "env"))
		m = get(Apps, "apps", i, "env").(map[string]interface{})
		m["id"] = get(Apps, "apps", i, "id")
		appsSlice = append(appsSlice, m)
	}

	return appsSlice
}

func get(m interface{}, path ...interface{}) interface{} {
	for _, p := range path {
		switch idx := p.(type) {
		case string:
			m = m.(map[string]interface{})[idx]
		case int:
			m = m.([]interface{})[idx]
		}
	}
	return m
}

func set(v interface{}, m interface{}, path ...interface{}) {
	for i, p := range path {
		last := i == len(path)-1
		switch idx := p.(type) {
		case string:
			if last {
				m.(map[string]interface{})[idx] = v
			} else {
				m = m.(map[string]interface{})[idx]
			}
		case int:
			if last {
				m.([]interface{})[idx] = v
			} else {
				m = m.([]interface{})[idx]
			}
		}
	}
}
