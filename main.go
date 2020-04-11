package main

import (
	"fmt"
	app "github.com/santinoncs/GoTicketHTTP/app"
	//"sync"
	"encoding/json"
	"net/http"
)

type Response struct {
	Status bool `json:"status"`
	Body   string `json:"body"`
}

// IncomingQuestion : here you tell us what IncomingQuestion is
type IncomingQuestion struct {
	Priority   int         `json:"priority"`
	Question   string      `json:"question"`
}


func main() {


	//jobQueue := make(chan app.Job, 100)

	var jobQueue = []chan app.Job {
		make(chan app.Job),
		make(chan app.Job),
	 }

	app.Start(jobQueue)
	

	//http.HandleFunc("/", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, jobQueue)
	})
	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)

}

// func response(status string, body string) {

// }

func handler(w http.ResponseWriter, r *http.Request,jobQueue []chan app.Job) {

	
	var response app.Response
	var content IncomingQuestion
	
	if r.URL.Path == "/api/post" {
		
		
		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				return
		}
		
		response = app.Post(content.Priority, content.Question,jobQueue)

	}
	

	fmt.Println("message respond is:", response.Message)
	//fmt.Println("Processed questions are:", st.GetProcessed())
	//fmt.Println("average_response_time in µs:", st.GetAverage())

	// }

	responseHTTP := Response{Status: response.Success, Body: response.Message}
	responseJSON, _ := json.Marshal(responseHTTP)
	fmt.Fprintf(w, "Response: %s\n", responseJSON)
}
