package viewmodels

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var Apps map[string]interface{}

func GetApps(w http.ResponseWriter, req *http.Request) []map[string]interface{} {
	// marathon ip also have to be imported to html form
	marathonIP := "http://10.198.161.42:8080/v2/apps"
	//marathonIP := "http://marathon-dev.mci01-uce.uce.corvisa.net:8080/v2/apps"

	err := req.ParseForm()
	if err != nil {
		fmt.Printf("ParseForm err : %s", err)
	}

	src := req.PostFormValue("source")
	tar := req.PostFormValue("target")
	// put some strict matching query for the rest call
	marathonURL := marathonIP + "?id=" + src
	res, err := http.Get(marathonURL)
	if err != nil {
		fmt.Printf("marathonURL err : %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("err : %s", err)
	}

	// returns nested map where the internal maps with keys will be of type map[string]interface{}
	// and the json arrays will be of type []interface which can be accessed through indices
	json.Unmarshal(body, &Apps)
	//fmt.Println(Apps)
	m1 := make(map[string]interface{})
	m2 := make(map[string]interface{})
	var appId []interface{}
	var appsSlice []map[string]interface{}
	var count int
	j := 0
	for i := range Apps["apps"].([]interface{}) {
		// Getting each app as map at the apps env field
		m1 = get(Apps, "apps", i, "env").(map[string]interface{})
		// what if env has also key as id
		appId = append(appId, get(Apps, "apps", i, "id"))

		// creating new map of each apps with id appended n few extra stuff
		// adding src n tar in 1st app map
		if count == 0 {
			m2["source"] = src
			m2["target"] = tar
			count++
		}
		//id is required to match with the src app
		m2["id"+strconv.Itoa(j)] = appId[j].(string)
		// newID is required to edit the app id
		m2["newID"] = appId[j]
		for key, value := range m1 {
			m2[key] = value
		}
		appsSlice = append(appsSlice, m2)
		j++
		m2 = make(map[string]interface{})
	}

	return appsSlice
}

// get set to nested maps
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
