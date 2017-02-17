package viewmodels

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func MigrateApps(w http.ResponseWriter, r *http.Request) interface{} {
	//marathonIP := "http://10.198.161.41:8080/v2/apps"
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("ParseForm err : %s", err)
	}
	src := r.PostFormValue("source")
	//tar := req.PostFormValue("target")
	//fmt.Printf(" src = %s ", src)
	if src == "" {
		fmt.Println("/apps has not redirected to /index, its we who is calling the /index")
		return nil
	}
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }

	// finding all the app ids as slices
	var idSlices []string
	var appsSlice []map[string]interface{}
	modifiedApp := make(map[string]interface{})
	i := 0
	for {
		// checking if id exist and getting all the ids in slices
		if id, ok := r.Form["id"+strconv.Itoa(i)]; ok {
			//fmt.Printf("id : %s", id) // getting 3 ids
			idSlices = append(idSlices, strings.Join(id, ""))
			// search for app attributes having id as 1st substring
			for k, v := range r.Form {
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
			appsSlice = append(appsSlice, modifiedApp)
			modifiedApp = make(map[string]interface{})
			i++
		} else {
			break
		}
	}
	fmt.Println(appsSlice)

	// modifiedApp := make(map[string]interface{})
	// i := 0
	// for k, v := range r.Form {
	// 	// getting apps based on id
	// 	idKey := r.Form["id"+strconv.Itoa(i)]
	// 	substr := k[len(idKey):len(k)]
	// 	i++
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }

	return "Successfully migrated the UC apps"
}

// func for original apps for the environment
