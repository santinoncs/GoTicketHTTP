package app

import (
	"fmt"
	"time"
	"sync"
)


// App : here you tell us what App is
type App struct {
	Status 
	jobQueue []chan Job
	workers int
	priority int
}

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
	mutex   			sync.Mutex
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

// NewApp : here you tell us what NewApp is
func NewApp(workers int) *App {
	arrQueue := make([]chan Job, 2)
	for i := range arrQueue {
		fmt.Println("the number is", i)
		arrQueue[i] = make(chan Job)
	}
	return &App{
		jobQueue: arrQueue,
		workers: workers,
	}
}

func newResponse(success bool, message string) *Response {
	return &Response{
		Success: success,
		Message: message,
	}
}

// Start : starting workers
func (a *App) Start() {


	for j := 1; j <= a.workers; j++ {
		a.Status.Workers ++
		go newWorker(j,a.jobQueue)
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
	//fmt.Println("The queue 0 is:", &jobQueue[0])
	//fmt.Printf("\nSlice: %T", jobQueue)
	//fmt.Printf("\nPointer: %d", len(jobQueue))

	for {
		select {
		case msg1 := <-jobQueue[0]:
			fmt.Println("listening on :", jobQueue[0])
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

	fmt.Println("accessing new job")

	responseChan1 := make(chan Response)

	j := Job{ID: priority, Question: question, ResponseChan: responseChan1}

	fmt.Println("creating job:", j)

	return j
}

// Post : escribo los jobs en jobs channel ya con los datos de prio y message
func (a *App) Post(priority int, question string) (Response) {

	start := time.Now()

	j := newJob(priority, question)

	fmt.Printf("El 0: %T", a.jobQueue[0])


	// aqui lanzo con go func el escribir en el channel de jobs

	go func() {

		if priority == 1 {
			a.jobQueue[0] <- j
		}
		if priority == 2 {
			a.jobQueue[1] <- j
		}

	}()

	channelListenR := j.ResponseChan

	select {
	case Response := <-channelListenR:
		t := time.Now()
		elapsed := t.Sub(start)
		a.Status.SetProcessed(elapsed)
		return Response
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
		res := newResponse(true, "error")
		return *res
	}

}

// SetProcessed : method SetProcessed
func (s *Status ) SetProcessed(e time.Duration) {
	s.mutex.Lock()
	s.TimeProcessed += e
	s.Processed ++
	s.mutex.Unlock()
}

// GetStatus : method GetStatus
func (s Status ) GetStatus() Status{
	s.AverageResponseTime = s.GetAverage()
	return s
}


// GetAverage : method GetAverage
func (s *Status ) GetAverage() int64{
	var microsperprocess int64
	micros := int64(s.TimeProcessed / time.Microsecond)
	if s.Processed > 0 {
		microsperprocess = micros / int64(s.Processed)
	} else {
		microsperprocess = 0
	}
	return microsperprocess
}
