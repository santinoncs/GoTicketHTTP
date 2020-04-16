package main

import (
	"fmt"
	app "github.com/santinoncs/GoTicketHTTP/app"
	//"sync"
	"encoding/json"
	"net/http"
	_ "log"
)

// IncomingQuestion : here you tell us what IncomingQuestion is
type IncomingQuestion struct {
	Priority   int         `json:"priority"`
	Question   string      `json:"question"`
}

var application *app.App

func main() {

	//st := app.NewStatus()
	//var mutex = &sync.Mutex{}


	application = app.NewApp(2)

	application.Start()
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)

}


func handler(w http.ResponseWriter, r *http.Request) {

	var value string
	var status bool
	
	var response app.Response
	var responseHTTP app.Response
	var responseHTTPStatus app.Status
	var content IncomingQuestion

	if r.URL.Path == "/api" {

		value = "The valid methods are: /api/post /api/status"
		status = true
		response := app.Response{Success: status, Message: value}
		responseJSON, _ := json.Marshal(response)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)
	
	}	

	
	if r.URL.Path == "/api/post" {
		
		
		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				return
		}
		
		response = application.Post(content.Priority, content.Question)
		responseHTTP = app.Response{Success: response.Success, Message: response.Message}
		responseJSON, _ := json.Marshal(responseHTTP)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	if r.URL.Path == "/api/status" {
		
		responseHTTPStatus = application.Status.GetStatus()
		responseJSON, _ := json.Marshal(responseHTTPStatus)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	
}
