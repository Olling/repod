package main

import (
	"os"
	"fmt"
	"flag"
	"sync"
	"time"
	"errors"
	"strings"
	"strconv"
	"os/exec"
	"net/http"
	"io/ioutil"
	"os/signal"
	"github.com/olling/logger"
	"github.com/olling/repod/conf"
)



var (
	ConfigurationPath = flag.String("configurationpath","/etc/repod/repod.conf","(Optional) The path to the configuration file")
)


func runScript(a conf.Action) (err error) {
	logger.Debug("Calling ActionChild: ", a)

	if a.Enabled == false {
		logger.Error("The ActionChild is disabled")
		return errors.New("The ActionChild is disabled and will not be executed")
	}

	cmd := exec.Command("/bin/bash", "-c", a.Command)
	output, err := cmd.CombinedOutput()

	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Debug("Script output:", string(output))
	return err
}


func httpDelete(w http.ResponseWriter, r *http.Request) {
	logger.Debug("DELETE url=" + r.URL.Path)
	fmt.Fprint(w,"Not Implemented")
}


func httpGet(w http.ResponseWriter, r *http.Request) {
	logger.Debug("GET url=" + r.URL.Path)

	if strings.Contains(r.URL.Path, "robots.txt") {
		logger.Debug("robots.txt ignored")
		fmt.Fprint(w,"Robots are not allowed")
		return
	}

	pathstat,err := os.Stat(conf.CurrentConfiguration.PathWork + r.URL.Path)
	if os.IsNotExist(err) {
		logger.Debug("Path: " +  conf.CurrentConfiguration.PathWork + r.URL.Path + " does not exist")
		fmt.Fprint(w,"The file does not exist")
		return
	}

	if ! pathstat.Mode().IsDir() {
                logger.Debug("Path: " +  conf.CurrentConfiguration.PathWork + r.URL.Path + " is not a directory")
		httpGetFile(w,r,conf.CurrentConfiguration.PathWork + r.URL.Path)
		return
	}

	httpGetDirectory(w,r)
}


func httpGetFile(w http.ResponseWriter, r *http.Request, path string) {
	for k,_ := range r.URL.Query() {
		if k == "edit" {
			logger.Debug("Edit Mode")
			return
		}
	}
	logger.Debug("View File Mode")

	a,err := conf.LoadActionFromPath(path)
	if err != nil {
		logger.Error("Could not read json file")
		http.Error(w, "Could not read json file" + err.Error(),500)
		return
	} 
	fmt.Fprint(w, a)
}


func httpGetDirectory(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(conf.CurrentConfiguration.PathWork + r.URL.Path)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintln(w,"<!DOCTYPE html>")
	fmt.Fprintln(w,"<html>")
	fmt.Fprintln(w,"<body>")

	for _, file := range files {
		prefix := ""
		if r.URL.Path != "/" {
			prefix = r.URL.Path
		}

		uri := prefix + "/" + file.Name()
		if file.IsDir() {
			fmt.Fprintln(w,"<a href=" + uri + ">" + file.Name() + "</a></br>")
		} else {
			fmt.Fprintln(w,file.Name() + " <a href=" + uri + ">View</a> <a href=" + uri + "?edit>Edit</a> <a href=" + uri + "?delete>Delete</a></br>")
		}
	}

	fmt.Fprintln(w,"</body>")
	fmt.Fprintln(w,"</html>")
}


func httpPost(w http.ResponseWriter, r *http.Request) {
	logger.Debug("POST url=" + r.URL.Path)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
                case "GET":
			httpGet(w,r)
		case "DELETE":
			httpDelete(w,r)	
                case "POST":
                        httpPost(w,r)
        }
}


func handleFavicon(w http.ResponseWriter, r *http.Request) {}
func handleStatus(w http.ResponseWriter, r *http.Request) {fmt.Fprint(w,"Running")}

func initializeInterupt(threadcount int, serverlistChannel chan string) (<-chan os.Signal, *sync.WaitGroup) {
	//To catch an interupt by the user
	interruptChannel := make(chan os.Signal, 1)

	//A channel to tell the threads to stop what they are doing
	killthreadsChannel := make(chan os.Signal, 1)

	//A wait group to wait for all of the threads to stop
	var interruptWaitGroup sync.WaitGroup
	interruptWaitGroup.Add(threadcount)

	//Catch user interrupt
	signal.Notify(interruptChannel, os.Interrupt)
	go func(){
	    for sig := range interruptChannel {
		logger.Error("Program was interrupted")

		//Empty the work list
		for server := range serverlistChannel {
			logger.Error("Removing the following server from the queue:", server)
		}

		//Tell the threads to interrupt
		for a := 1; a <= threadcount; a++ {
			killthreadsChannel <-sig
		}

		//Wait for all the groups to stop
		interruptWaitGroup.Wait()

		//Exit
		os.Exit(1)
	    }
	}()
	return killthreadsChannel, &interruptWaitGroup
}


func main() {
	logger.Initialize()

	conf.WriteJsonFile(time.Now(), "/tmp/test/time")
	logger.Info(time.Now())
	logger.Info("starting test")
	logger.Info(conf.LoadChannelsFromPath("/tmp/test"))

	//Parse the flags/arguments to the properties
	flag.Parse()
	conf.InitializeConfiguration(*ConfigurationPath)
	logger.SetDebugState(conf.CurrentConfiguration.Debug)

	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/favicon.ico", handleFavicon)
        http.HandleFunc("/", httpHandler)

	logger.Info("Listening on port: " + strconv.Itoa(conf.CurrentConfiguration.TlsPort) + " (https)")
	tlserr := http.ListenAndServeTLS(":" + strconv.Itoa(conf.CurrentConfiguration.TlsPort), conf.CurrentConfiguration.TlsCert, conf.CurrentConfiguration.TlsKey, nil)

	if tlserr != nil {
		logger.Error("Error starting TLS: ",tlserr)
	}

	logger.Error("Error happend while serving website")
}
