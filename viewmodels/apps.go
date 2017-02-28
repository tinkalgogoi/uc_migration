package viewmodels

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var Apps map[string]interface{}
var AppsTarget map[string]interface{}

var srcAppsSlice []map[string]interface{}
var tarAppsSlice []map[string]interface{}

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

	// getting app json from src
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

	// getting app json from tar
	marathonURL = marathonIP + "?id=" + tar
	res, err = http.Get(marathonURL)
	if err != nil {
		fmt.Printf("marathonURL err : %s", err)
	}

	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("err : %s", err)
	}
	json.Unmarshal(body, &AppsTarget)

	srcAppsSlice := getAppsEnv(src, tar, "src", Apps)
	fmt.Println("srcAppsSlice", srcAppsSlice)
	tarAppsSlice := getAppsEnv(src, tar, "tar", AppsTarget)
	fmt.Println("tarAppsSlice", tarAppsSlice)

	appsSlice := populateEnv(srcAppsSlice, tarAppsSlice)

	return appsSlice

}

// populating target environment available env's to source environment env's
func populateEnv(srcSlice []map[string]interface{}, tarSlice []map[string]interface{}) []map[string]interface{} {
	for i := range srcSlice {
		for j := range tarSlice { // refactor the code where once checked, it shouldn't again check the app

			if srcSlice[i]["appName"] == tarSlice[j]["appName"] {
				for key, value := range tarSlice[j] {
					match, _ := regexp.MatchString("^id([0-9]+)$", key)
					if match {
						// not adding the id.. key of target application in modified app as we need the src id (tar id will be in newID key)
						continue
					} else {
						srcSlice[i][key] = value
					}
				}

			}
		}

	}
	return srcSlice
}

func appNamebetween(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

// get the environments of a app with id embedded in each app map
func getAppsEnv(src string, tar string, appType string, apps map[string]interface{}) []map[string]interface{} {
	m1 := make(map[string]interface{})
	m2 := make(map[string]interface{})
	var appId []interface{}
	var appsSlice []map[string]interface{}
	var count int
	j := 0
	for i := range apps["apps"].([]interface{}) {
		// Getting each app as map at the apps env field
		m1 = get(apps, "apps", i, "env").(map[string]interface{})
		// what if env has also key as id
		appId = append(appId, get(apps, "apps", i, "id"))

		// app name of the application between the application environment type ex. /dev/uc and /api substring
		if appType == "src" {
			m2["appName"] = appNamebetween(appId[j].(string), src, "/api")
		} else if appType == "tar" {
			m2["appName"] = appNamebetween(appId[j].(string), tar, "/api")
		}

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
