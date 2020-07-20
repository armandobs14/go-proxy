package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var portMapping = map[string]string{
	"8081": "9999",
	"8082": "9999",
}

func fowardRequest(w http.ResponseWriter, r *http.Request) {
	var portForward string
	if i := strings.Index(r.Host, ":"); i != -1 {
		port := r.Host[i+1:]
		portForward = portMapping[port]
	}

	urlString := fmt.Sprintf("http://localhost:%s", portForward)
	urlPath, _ := url.Parse(urlString)
	httputil.NewSingleHostReverseProxy(urlPath).ServeHTTP(w, r)
}

func main() {
	for key, value := range portMapping {
		fmt.Println("Ingress:", key, "Forward:", value)
		go func(key string, value string) {
			mux := http.NewServeMux()
			mux.HandleFunc("/", fowardRequest)
			err := http.ListenAndServe(fmt.Sprintf(":%s", key), mux)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		}(key, value)

	}
	select {}
}
