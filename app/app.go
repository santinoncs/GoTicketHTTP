package app

import (
	"fmt"
	"time"
	//"sync"
)



// Response : here you tell us what Response is
type Response struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
}


// Status : here you tell us what Status is
type Status struct {
	Workers             int
	Processed           int
	TimeProcessed       time.Duration
	AverageResponseTime int64
}

// Job : here you tell us what Job is
type Job struct {
	ID           int
	Question     string
	ResponseChan chan Response
}

// NewStatus : Constructor of status struct
func NewStatus() *Status {
	return &Status{}
}

func newResponse(success bool, message string) *Response {
	return &Response{
		Success: success,
		Message: message,
	}
}

// Start : starting workers
func Start(jobQueue []chan Job,st *Status) {

	fmt.Println("Starting the workers")

	numWorkers := 2

	for j := 1; j <= numWorkers; j++ {
		st.Workers ++
		go newWorker(j,jobQueue)
	}
}

func (j Job) process() Response {

	var mess string	
	
	if j.ID == 1 {
		mess = "This is prio 1"
	} 
	if j.ID == 2 {
		mess = "This is prio 2"
	}
	
	
	res := newResponse ( true, mess )
	return *res

}

func newWorker(j int,jobQueue []chan Job) {

	fmt.Println("worker started:", j)

	for {
		select {
		case msg1 := <-jobQueue[0]:
			//time.Sleep(1 * time.Second)
			msg1.ResponseChan <- msg1.process()
			close(msg1.ResponseChan)
		case msg2 := <-jobQueue[1]:
			//time.Sleep(1 * time.Second)
			msg2.ResponseChan <- msg2.process()
			close(msg2.ResponseChan)
		}
		
	}
	

}

func newJob(priority int, question string) Job {

	responseChan1 := make(chan Response)

	j := Job{ID: priority, Question: question, ResponseChan: responseChan1}

	return j
}

// Post : escribo los jobs en jobs channel ya con los datos de prio y message
func Post(priority int, question string, jobQueue []chan Job,st *Status) (Response) {

	start := time.Now()

	fmt.Println("Entering Post...")

	j := newJob(priority, question)


	// aqui lanzo con go func el escribir en el channel de jobs

	go func() {

		if priority == 1 {
			jobQueue[0] <- j
		}
		if priority == 2 {
			fmt.Println("sending to jobqueue2")
			jobQueue[1] <- j
		}

	}()

	channelListenR := j.ResponseChan

	select {
	case Response := <-channelListenR:
		t := time.Now()
		elapsed := t.Sub(start)
		//mutex.Lock()
		st.TimeProcessed += elapsed
		st.Processed ++
		//mutex.Unlock()
		return Response
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
		res := newResponse(true, "error")
		return *res
	}

}

func (st *Status ) GetProcessed() int{
	 return st.Processed
}

func (st *Status ) GetWorkers() int{
	 return st.Workers
}

func (st *Status ) GetAverage() int64{
	micros := int64(st.TimeProcessed / time.Microsecond)
	return micros
}
