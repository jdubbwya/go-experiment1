package server

import (
	"os"
	"time"
)

func notifyShutdown(){
	close(Monitor.Channel)
}

func drainConnectionsThenShutdown(){
	Monitor.WaitGroup.Wait()
	time.Sleep(time.Second)
	os.Exit(0)
}