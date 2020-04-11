package app

import (
	"fmt"
	"time"
	//"sync"
)

// var jobchan1 chan Job
// var jobchan2 chan Job



// type jobchan struct {
// 	jobchan chan Job
// }



// Response : here you tell us what Response is
type Response struct {
	Success bool
	Message string
}

// Status : here you tell us what Status is
type Status struct {
	workers   int
	processed int
	timeProcessed   time.Duration
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
func Start(jobQueue []chan Job) {

	fmt.Println("Starting the workers")

	numWorkers := 2

	for j := 1; j <= numWorkers; j++ {
		go newWorker(j,jobQueue)
	}
}

func (j Job) process() Response {

	res := newResponse ( true, "bye" )
	return *res

}

func newWorker(j int,jobQueue []chan Job) {

	fmt.Println("worker started:", j)

	for {
		select {
		case msg1 := <-jobQueue[0]:
			time.Sleep(1 * time.Second)
			fmt.Println("Receiving from channel jobchan1 and writing to chan:", msg1.ResponseChan)
			msg1.ResponseChan <- msg1.process()
			close(msg1.ResponseChan)
		case msg2 := <-jobQueue[1]:
			time.Sleep(4 * time.Second)
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
func Post(priority int, question string, jobQueue []chan Job) (Response) {

	// start := time.Now()

	fmt.Println("Entering Post...")

	j := newJob(priority, question)

	// jobchan1 = make(chan Job)
	// jobchan2 = make(chan Job)

	// aqui lanzo con go func el escribir en el channel de jobs

	go func() {

		if priority == 1 {
			fmt.Println("Sending job to channel jobchan1", j)
			jobQueue[0] <- j
		}
		if priority == 2 {
			fmt.Println("Sending jo to channel jobchan2", j)
			jobQueue[1] <- j
		}

	}()

	channelListenR := j.ResponseChan

	select {
	case Response := <-channelListenR:
		// t := time.Now()
		// elapsed := t.Sub(start)
		// mutex.Lock()
		// st.timeProcessed += elapsed
		// st.processed ++
		// mutex.Unlock()
		return Response
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
		res := newResponse(true, "error")
		return *res
	}

}

func (st *Status ) GetProcessed() int{
	 return st.processed
}

func (st *Status ) GetWorkers() int{
	 return st.workers
}

func (st *Status ) GetAverage() int64{
	micros := int64(st.timeProcessed / time.Microsecond)
	return micros
}
