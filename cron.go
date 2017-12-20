package main

import (
	"github.com/robfig/cron"
	"github.com/olling/logger"
)

var (
	CurrentCron = cron.New()
)

func initializeCron() {
	logger.Debug("Initializing Cron")
	CurrentCron.Start()
	logger.Debug("Started cron")
}


func (action *Action) StartCron () {
	CurrentCron.AddFunc(action.Cron, action.cronTick)
	//action.cronIndex = len(CurrentCron.Entries()) -1
}


//func (action *Action) StopCron() {
//	logger.Debug("I AM", action.cronIndex)
//}


func (action Action) cronTick () {
	logger.Debug(action)
}



//func cronTick (action Action) {
//	logger.Debug(action)
//
//	//For each child tick on them
//
//	//execute Serverinfo
//}
//func cronTick (channel Channel) {
//	logger.Debug(channel)
//
//	//For each child tick on them
//
//	//execute Serverinfo
//}
