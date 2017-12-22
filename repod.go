package main

import (
	//"os"
	//"sync"
	//"time"
	"strconv"
	"net/http"
	//"os/signal"
	"github.com/olling/logger"
	//"github.com/olling/repod/conf"
)


//func initializeInterupt(threadcount int, serverlistCannel chan string) (<-chan os.Signal, *sync.WaitGroup) {
//	//To catch an interupt by the user
//	interruptChannel := make(chan os.Signal, 1)
//
//	//A channel to tell the threads to stop what they are doing
//	killthreadsChannel := make(chan os.Signal, 1)
//
//	//A wait group to wait for all of the threads to stop
//	var interruptWaitGroup sync.WaitGroup
//	interruptWaitGroup.Add(threadcount)
//
//	//Catch user interrupt
//	signal.Notify(interruptChannel, os.Interrupt)
//	go func(){
//	    for sig := range interruptChannel {
//		logger.Error("Program was interrupted")
//
//		//Empty the work list
//		for server := range serverlistChannel {
//			logger.Error("Removing the following server from the queue:", server)
//		}
//
//		//Tell the threads to interrupt
//		for a := 1; a <= threadcount; a++ {
//			killthreadsChannel <-sig
//		}
//
//		//Wait for all the groups to stop
//		interruptWaitGroup.Wait()
//
//		//Exit
//		os.Exit(1)
//	    }
//	}()
//	return killthreadsChannel, &interruptWaitGroup
//}


func main() {
	logger.Initialize()

	//WriteJsonFile(time.Now(), "/tmp/test/time")
	//logger.Info(time.Now())
	//logger.Info("starting test")
	//logger.Info(LoadChannelsFromPath("/tmp/test"))

	//Parse the flags/arguments to the properties
	InitializeConfiguration()
	logger.SetDebugState(CurrentConfiguration.Debug)

	initializeCron()

	//	var test Action
	//	test.Cron = "@every 3s"
	//	test.StartCron()
	//
	//	var test2 Action
	//	test2.Cron ="@every 4s"
	//	test2.StartCron()

	initializeApi()
	initializeWebinterface()

	logger.Debug("here")
	actions,_ := LoadActionsFromPath(CurrentConfiguration.PathWork)
	StartCrons(actions)

	if CurrentConfiguration.HttpTlsPort != 0 {
		logger.Info("Listening on port: " + strconv.Itoa(CurrentConfiguration.HttpTlsPort) + " (https)")
		go http.ListenAndServeTLS(":" + strconv.Itoa(CurrentConfiguration.HttpTlsPort), CurrentConfiguration.HttpTlsCert, CurrentConfiguration.HttpTlsKey, nil)
	}


	if CurrentConfiguration.HttpPort != 0 {
		logger.Info("Listening on port: " + strconv.Itoa(CurrentConfiguration.HttpPort) + " (http)")
		http.ListenAndServe(":" + strconv.Itoa(CurrentConfiguration.HttpPort),nil)
	}

	logger.Error("Error happend while serving website")
}
