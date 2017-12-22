package main

import (
	"github.com/robfig/cron"
	"github.com/olling/logger"
	"os/exec"
)

var (
	CurrentCron = cron.New()
)

func initializeCron() {
	logger.Debug("Initializing Cron")
	CurrentCron.Start()
	logger.Debug("Started cron")
}

func StartCrons (actions []Action) {
	logger.Debug("StartCron on array")
	for _,action := range actions {
		action.StartCron()
		if len(action.children) > 0 {
			StartCrons(action.children)
		}
	}
}

func (action *Action) StartCron () {
	logger.Info("Starting action:", action.FriendlyName)
	CurrentCron.AddFunc(action.Cron, action.cronTick)
}


func (action *Action) Execute () {
	if action.Enabled == false {
		logger.Debug("Ignoring disabled action " + action.Name + "(" + action.FriendlyName + ")")
		return
	}
	logger.Info("Executing action" + action.Name + "(" + action.FriendlyName + ")")

	cmd := exec.Command("/bin/bash", "-c", action.Command)
	output, err := cmd.CombinedOutput()

	if err != nil {
		logger.Error(err)
	}

	logger.Info("Executed action" + action.Name + "(" + action.FriendlyName + ")")
	logger.Debug("Script output:", string(output))

	for _,child := range action.children {
		logger.Debug(child)
		logger.Debug("Executing children")
		child.Execute()
	}
}

func (action Action) cronTick () {
	logger.Debug(action)
	action.Execute()
}
