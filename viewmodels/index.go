package viewmodels

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var newAppsSlice []map[string]interface{}
var oldAppsSlice map[string]interface{}
var src string
var tar string

func MigrateApps(w http.ResponseWriter, r *http.Request) interface{} {
	//marathonIP := "http://10.198.161.41:8080/v2/apps"
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("ParseForm err : %s", err)
	}
	src = r.PostFormValue("source")
	//tar = req.PostFormValue("target")

	if src == "" {
		// the /apps has not redirected to /index, its we who is calling the /index
		return nil
	}

	// finding all the app ids as slices
	var idSlices []string

	modifiedApp := make(map[string]interface{})
	i := 0
	for {
		// checking if id exist and getting all the ids in slices
		if id, ok := r.Form["id"+strconv.Itoa(i)]; ok {
			//fmt.Printf("id : %s", id) // getting 3 ids
			idSlices = append(idSlices, strings.Join(id, ""))
			// search for app attributes having id as 1st substring
			for k, v := range r.Form { // code can be refactored, right now matching every key for each id
				if len(k) >= len(idSlices[i]) { // gives error while matching if key length is less than id
					// matching id with map keys
					if idSlices[i] == k[:len(idSlices[i])] {
						// getting other part of env attribute string which doesn't match with id
						modifiedApp[k[len(idSlices[i]):len(k)]] = v
					}
				}
			}
			modifiedApp["id"] = idSlices[i]
			// got the  map of app and putting it in slices
			newAppsSlice = append(newAppsSlice, modifiedApp)
			modifiedApp = make(map[string]interface{})
			i++
		} else {
			break
		}
	}

	getOldApps()
	createNewApp()
	fmt.Println(oldAppsSlice)

	//fmt.Println(newAppsSlice)

	return "Successfully migrated the UC apps"
}

// func for original apps for the environment
func getOldApps() map[string]interface{} {
	marathonIP := "http://10.198.161.41:8080/v2/apps"
	marathonURL := marathonIP + "?id=" + src
	res, err := http.Get(marathonURL)
	if err != nil {
		fmt.Printf("marathonURL err : %s", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("err : %s", err)
	}

	json.Unmarshal(body, &oldAppsSlice)
	return oldAppsSlice
}

// Creating new app with new environment configuration
// few additional attributes need to be removed before posting the app
func createNewApp() map[string]interface{} {
	//m1 := make(map[string]interface{})
	var appId interface{}
	for i := range oldAppsSlice["apps"].([]interface{}) { // iterating all apps because need to remove few attributes ex. versions etc

		// what if env has also key as id
		appId = get(oldAppsSlice, "apps", i, "id")

		// Matching the apps and merging the new environment
		for j := range newAppsSlice { // code can be refactored as we are matching each app with modified app
			// matching the old ids of both app
			if newAppsSlice[j]["id"] == appId {
				//setting the new id to the particular app
				//fmt.Printf("ids are : %s", newAppsSlice[j]["id"])
				set(newAppsSlice[j]["newID"], oldAppsSlice, "apps", i, "id")
				// deleting unwanted id and oldID values from app env
				delete(newAppsSlice[j], "id")
				delete(newAppsSlice[j], "newID")
				// setting my new evrionment against old environment
				set(newAppsSlice[j], oldAppsSlice, "apps", i, "env")
			}
		}

		// remove attributes causing problem on posting maarathon app
	}

	return oldAppsSlice // this old map is now populated with new values
}

// post the new app to marathon
