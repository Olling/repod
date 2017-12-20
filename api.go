package main

import (
	"fmt"
	"net/http"
	"github.com/olling/logger"
)

func httpApiChannels(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Api: Channels called")
        switch r.Method {
                case "GET":
			httpApiGet(w,r)
		case "DELETE":
			httpApiDelete(w,r)
                case "POST":
                        httpApiPost(w,r)
        }
}

func httpApiGet(w http.ResponseWriter, r *http.Request) {
	logger.Debug("API called with GET", r.URL.Path)
}
func httpApiPost(w http.ResponseWriter, r *http.Request) {
	logger.Debug("API called with POST", r.URL.Path)
}
func httpApiDelete(w http.ResponseWriter, r *http.Request) {
	logger.Debug("API called with Delete", r.URL.Path)
}

func httpApiHandler(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
                case "/api/channels":
			httpApiChannels(w,r)
        }
}

func handleApiFavicon(w http.ResponseWriter, r *http.Request) {}
func handleApiStatus(w http.ResponseWriter, r *http.Request) {fmt.Fprint(w,"Running")}

func initializeApi() {
	http.HandleFunc("/api/status", handleApiStatus)
	http.HandleFunc("/api/favicon.ico", handleApiFavicon)
	http.HandleFunc("/api", httpApiHandler)
}
