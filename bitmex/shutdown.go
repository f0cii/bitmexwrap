package bitmex

import log "github.com/sirupsen/logrus"

// Shutdown to monitor and shut down routines package specific
type Shutdown struct {
	c            chan int
	routineCount int
	finishC      chan int
}

// NewRoutineManagement returns an new initial routine management system
func NewRoutineManagement() *Shutdown {
	return &Shutdown{
		c:       make(chan int, 1),
		finishC: make(chan int, 1),
	}
}

// AddRoutine adds a routine to the monitor and returns a channel
func (r *Shutdown) addRoutine() chan int {
	log.Println("Bitmex Websocket: Routine added to monitor")
	r.routineCount++
	return r.c
}

// RoutineShutdown sends a message to the finisher channel
func (r *Shutdown) routineShutdown() {
	log.Println("Bitmex Websocket: Routine is shutting down")
	r.finishC <- 1
}

// SignalShutdown signals a shutdown across routines
func (r *Shutdown) SignalShutdown() {
	log.Println("Bitmex Websocket: Shutdown signal sending..")
	for i := 0; i < r.routineCount; i++ {
		log.Printf("Bitmex Websocket: Shutdown signal sent to routine %d", i+1)
		r.c <- 1
	}

	for {
		<-r.finishC
		r.routineCount--
		if r.routineCount <= 0 {
			close(r.c)
			close(r.finishC)
			log.Println("Bitmex Websocket: All routines stopped")
			return
		}
	}
}
