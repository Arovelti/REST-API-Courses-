package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type courseInfo struct {
	Title string `json: "Title"`
}

var courses map[string]courseInfo

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API!")
}

func allcourses(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "List of all courses")
	kv := r.URL.Query()

	for k, v := range kv {
		fmt.Println(k, v)
	}

	//Example: check for the "country" key
	//if val, ok := kv["country"]; ok {
	//	fmt.Println(val[0])
	//}

	json.NewEncoder(w).Encode(courses)
}

func course(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//fmt.Fprintf(w, "Detail for course"+params["courseid"])
	//fmt.Fprintf(w, "\n")
	//fmt.Fprintf(w, r.Method)

	if r.Method == "GET" {
		if _, ok := courses[params["courseid"]]; ok {
			json.NewEncoder(w).Encode(courses[params{"courseid"}])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
		}
	}

	if r.Method == "DELETE" {
		if _, ok := courses[params["courseid"]]; ok {
			delete(courses, params["courseid"])
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
		}
	}

	if r.Header.Get("Content-type") == "aplication/json" {
		// for creating a new course
		if r.Method == "POST" {
			//read the string sent to the service
			var newCourse courseInfo
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				//convert json to object
				json.Unmarshal(reqBody, &newCourse)

				if newCourse.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity) //Необработанный объект
					w.Write([]byte("422 - Please supply course " + "information " + "in JSON format"))
					return
				}
				//check if course exists / add only when it doesn't exist
				if _, ok := courses[params["courseid"]]; !ok {
					courses[params["courseid"]] = newCourse
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " + params["courseid"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - Duplicate course id"))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " + "course information in JSON format"))
			}
		}
		// creating or updating
		if r.Method == "PUT" {
			var newCourse courseInfo
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newCourse)

				if newCourse.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply course " + "information " + "in JSON format"))
					return
				}
				// check if couse exists / add only if couse does't exist
				if _, ok := courses[params["courseid"]]; !ok {
					courses[params["courseid"]] = newCourse
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " + params["courseid"]))
				} else {
					//update course
					courses[params["courseid"]] = newCourse
					w.WriteHeader(http.StatusNoContent)
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " + "course information" + "in JSON format"))
			}
		}
	}
}

func main() {
	//instantiate courses
	courses = make(map[string]courseInfo)

	router := mux.NewRoater()
	router.HandleFunc("/api/v1", home)

	router.HandleFunc("api/v1/courses", allcourses)
	router.HandleFunc("api/v1/courses/{courseid}",
		course).Methods(
		"GET", "POST", "PUT", "DELETE")

	fmt.Println("Listen and serve at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
