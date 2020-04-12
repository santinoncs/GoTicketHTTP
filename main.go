package main

import (
	"fmt"
	app "github.com/santinoncs/GoTicketHTTP/app"
	"sync"
	"encoding/json"
	"net/http"
	_ "log"
)

// IncomingQuestion : here you tell us what IncomingQuestion is
type IncomingQuestion struct {
	Priority   int         `json:"priority"`
	Question   string      `json:"question"`
}

func main() {

	st := app.NewStatus()
	var mutex = &sync.Mutex{}

	var jobQueue = []chan app.Job {
		make(chan app.Job, 100),
		make(chan app.Job, 100),
	 }

	app.Start(jobQueue,st)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, jobQueue, st,mutex)
	})
	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)

}


func handler(w http.ResponseWriter, r *http.Request,jobQueue []chan app.Job, st *app.Status,mutex *sync.Mutex) {

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
		
		response = app.Post(content.Priority, content.Question,jobQueue,st,mutex)
		responseHTTP = app.Response{Success: response.Success, Message: response.Message}
		responseJSON, _ := json.Marshal(responseHTTP)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	if r.URL.Path == "/api/status" {
		
		NumberOfWorkers      := st.GetWorkers()
		NumberOfProcesses    := st.GetProcessed()
		AverageResponseTime  := st.GetAverage()

		responseHTTPStatus = app.Status{Workers: NumberOfWorkers, Processed: NumberOfProcesses, AverageResponseTime: AverageResponseTime}
		responseJSON, _ := json.Marshal(responseHTTPStatus)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	
}
