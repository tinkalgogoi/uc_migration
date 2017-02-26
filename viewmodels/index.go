package viewmodels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var newAppsSlice []map[string]interface{}
var oldAppsSlice map[string]interface{}
var src string
var tar string
var marathonIP string

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

	// getting all the app ids from the Form
	var keyIDs []string
	for kid := range r.Form {
		if strings.HasPrefix(kid, "id") {
			keyIDs = append(keyIDs, kid)
		}
	}

	for i, keyID := range keyIDs {

		id := r.Form[keyID]
		//fmt.Printf("id : %s", id)
		idSlices = append(idSlices, strings.Join(id, ""))
		// search for app attributes having id as 1st substring
		for k, v := range r.Form { // code can be refactored, right now matching every key for each id
			if len(k) >= len(idSlices[i]) { // gives error while matching if key length is less than id
				// matching id with map keys
				if idSlices[i] == k[:len(idSlices[i])] {
					// getting other part of env attribute string which doesn't match with id
					modifiedApp[k[len(idSlices[i]):len(k)]] = strings.Join(v, "") // converting v slice to string or else it will be marathon error since it will look like json array []
				}
			}
		}
		modifiedApp["id"] = idSlices[i]
		// got the  map of app and putting it in slices
		newAppsSlice = append(newAppsSlice, modifiedApp)
		modifiedApp = make(map[string]interface{})

	}

	fmt.Println("newAppsSlice :")
	fmt.Println(newAppsSlice)

	getOldApps()
	createNewApp()
	putApp()

	return "Successfully migrated the UC apps"
}

// func for original apps for the environment
func getOldApps() map[string]interface{} {
	marathonIP = "http://10.198.161.42:8080/v2/apps"
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
				//setting the new id to the particular app id field of oldAppsSlice
				set(newAppsSlice[j]["newID"], oldAppsSlice, "apps", i, "id")
				// deleting unwanted id and oldID values from app env of New app before merging to oldAppsSlice
				delete(newAppsSlice[j], "id")
				delete(newAppsSlice[j], "newID")

				// setting my new evrionment against old environment
				set(newAppsSlice[j], oldAppsSlice, "apps", i, "env")
			}
		}

		// remove attributes causing problem on posting maarathon app
		deleteKey(oldAppsSlice, "apps", i, "version")
	}
	fmt.Println("oldAppsSlice with new env values")
	fmt.Println(oldAppsSlice)
	return oldAppsSlice // this old map is now populated with new values
}

// post the new app to marathon
func putApp() interface{} {
	// use groups url as we might have groups inside our app
	marathonURL := "http://10.198.161.42:8080/v2/groups"
	appJson, err := json.Marshal(oldAppsSlice)
	if err != nil {
		fmt.Printf("err : %s", err)
	}
	fmt.Println("The json :")
	fmt.Println(string(appJson))
	//resp, err := http.Post(marathonURL, "application/json", bytes.NewBuffer(appJson))
	client := &http.Client{}
	request, err := http.NewRequest("PUT", marathonURL, bytes.NewBuffer(appJson))
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("err : %s", err)
	}
	fmt.Println("response Status:", response.Status)
	return nil
}

func deleteKey(m interface{}, path ...interface{}) {
	for i, p := range path {
		last := i == len(path)-1
		switch idx := p.(type) {
		case string:
			if last {
				delete(m.(map[string]interface{}), "version")
			} else {
				m = m.(map[string]interface{})[idx]
			}
		// zzzzzzzzzzzz
		case int:
			if last {
				// need to find a way to delete a type []interface{} which holds other maps (these are json arrays)
				// delete needs a map as parameter to delete, not []interface{} type
				fmt.Println("zzzzzzzzzzzzzz")
			} else {
				m = m.([]interface{})[idx]
			}
		}
	}
}
