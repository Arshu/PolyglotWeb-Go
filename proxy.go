package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	realBackend := os.Getenv("ARSHU_WEB")
	if realBackend == "" {
		realBackend = "http://localhost:5000"
	}
	realBackendUrl, err := url.Parse(realBackend)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to connect to backend", http.StatusBadGateway)
		return
	}

	clonedR := r.Clone(context.TODO())
	clonedR.URL.Host = realBackendUrl.Host
	clonedR.URL.Scheme = realBackendUrl.Scheme
	clonedR.RequestURI = ""
	response, err := http.DefaultClient.Do(clonedR)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to connect to backend", http.StatusBadGateway)
		return
	}

	for header, values := range response.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(response.StatusCode)

	io.Copy(w, response.Body)
}
