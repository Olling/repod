package main

import (
	"os"
	"fmt"
	//"time"
	"strings"
	//"strconv"
	"net/http"
//	"io/ioutil"
	"github.com/olling/logger"
	//"github.com/olling/repod/conf"
)


func httpDelete(w http.ResponseWriter, r *http.Request) {
	logger.Debug("DELETE url=" + r.URL.Path)
	fmt.Fprint(w,"Not Implemented")
}


func httpGet(w http.ResponseWriter, r *http.Request) {
	logger.Debug("GET url=" + r.URL.Path)

	if strings.Contains(r.URL.Path, "robots.txt") {
		logger.Debug("robots.txt ignored for now")
		fmt.Fprint(w,"Robots are not allowed - yet")
		return
	}

	pathstat,err := os.Stat(CurrentConfiguration.PathWork + r.URL.Path)
	if os.IsNotExist(err) {
		logger.Debug("Path: " +  CurrentConfiguration.PathWork + r.URL.Path + " does not exist")
		fmt.Fprint(w,"The file does not exist")
		return
	}

	if ! pathstat.Mode().IsDir() {
                logger.Debug("Path: " +  CurrentConfiguration.PathWork + r.URL.Path + " is not a directory")
		httpGetFile(w,r,CurrentConfiguration.PathWork + r.URL.Path)
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

	a,err := LoadActionFromPath(path)
	if err != nil {
		logger.Error("Could not read json file")
		http.Error(w, "Could not read json file" + err.Error(),500)
		return
	}
	fmt.Fprint(w, a)
}


func httpGetDirectory(w http.ResponseWriter, r *http.Request) {
	prefix := ""
	if r.URL.Path != "/" {
		prefix = r.URL.Path
	}

	fmt.Fprintln(w,"<!DOCTYPE html>")
	fmt.Fprintln(w,"<html>")
	fmt.Fprintln(w,"<body>")

	path := CurrentConfiguration.PathWork + prefix
	fileinfo,err := os.Stat(path)

	if err != nil {
		logger.Error("Could not stat file", path, err)
	}

	if fileinfo.IsDir() {
			actions,_ := LoadActionsFromPath(path)
			logger.Debug(actions)
			for _, action := range actions {
				uri := prefix + "/"  + action.fileName


				dirName := action.GetDirName()
				var children string
				if dirName == "" || len(action.children) == 0 {
					children = ""
				} else {
					children = " <a href=" + dirName + ">Children</a>"
				}
				fmt.Fprintln(w,action.FriendlyName + " <a href=" + uri + ">View</a>" + children + " <a href=" + uri + "?edit>Edit</a> <a href=" + uri + "?delete>Delete</a></br>")
			}
		} else {
			action,_ := LoadActionFromPath(path)
			uri := prefix + "/"  + action.fileName
			fmt.Fprintln(w,action.FriendlyName + " <a href=" + uri + ">View</a> <a href=" + uri + "?edit>Edit</a> <a href=" + uri + "?delete>Delete</a></br>")
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


func initializeWebinterface() {
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/favicon.ico", handleFavicon)
        http.HandleFunc("/", httpHandler)
//
//	if CurrentConfiguration.HttpTlsPort != 0 {
//		logger.Info("Listening on port: " + strconv.Itoa(CurrentConfiguration.HttpTlsPort) + " (https)")
//		go http.ListenAndServeTLS(":" + strconv.Itoa(CurrentConfiguration.HttpTlsPort), CurrentConfiguration.HttpTlsCert, CurrentConfiguration.HttpTlsKey, nil)
//	}
//
//
//	if CurrentConfiguration.HttpPort != 0 {
//		logger.Info("Listening on port: " + strconv.Itoa(CurrentConfiguration.HttpPort) + " (http)")
//		go http.ListenAndServe(":" + strconv.Itoa(CurrentConfiguration.HttpPort),nil)
//	}
//	logger.Error("Error happend while serving the webinterface")
}
