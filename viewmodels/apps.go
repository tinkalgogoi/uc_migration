package viewmodels

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var Apps map[string]interface{}

func GetApps(w http.ResponseWriter, req *http.Request) (Apps map[string]interface{}) {
	marathonIP := "http://10.198.161.41:8080/v2/groups/"

	err := req.ParseForm()
	if err != nil {
		fmt.Printf("ParseForm err : %s", err)
	}

	src := req.PostFormValue("source")
	//tar := req.PostFormValue("target")
	marathonURL := marathonIP + src
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
	fmt.Println(get(Apps, "apps", 0, "id"))
	fmt.Print("Environment : ")
	fmt.Println(get(Apps, "apps", 0, "env"))
	//fmt.Printf("Results: %v\n", body)

	return Apps
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
