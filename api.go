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
			httpGet(w,r)
		case "DELETE":
			httpDelete(w,r)
                case "POST":
                        httpPost(w,r)
        }
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
