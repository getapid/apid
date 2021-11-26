package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func any(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			w.Header().Add(name, h)
		}
	}

	if code, err := strconv.ParseUint(w.Header().Get("X-ECHO-STATUSCODE"), 10, 64); err == nil {
		w.WriteHeader(int(code))
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("cant read request body")
		return
	}
	w.Write(body)
}

func main() {
	http.HandleFunc("/", any)
	http.ListenAndServe(":8080", nil)
}
