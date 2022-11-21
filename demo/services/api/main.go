package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var serviceToken string

func readToken() {
	b, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		panic(err)
	}
	serviceToken = string(b)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NEW REQUEST")
	// Make a HTTP request to service2
	serviceConnstring := os.Getenv("STORAGE_HUB_SVC")
	if len(serviceConnstring) == 0 {
		panic("STORAGE_HUB_SVC expected")
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", serviceConnstring, nil)
	if err != nil {
		panic(err)
	}
	// Identity self to service 2 using service account token
	req.Header.Add("Bearer", serviceToken)
	// Add podIP on header
	req.Header.Add("X-Real-IP", os.Getenv("MY_POD_IP"))
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		io.WriteString(w, string(body))
	}
	fmt.Println("RESPONSE OK")

}

func main() {
	fmt.Println("start...")

	// Read the token at startup
	readToken()
	http.HandleFunc("/", handleIndex)
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		panic("LISTEN_ADDR expected")
	}
	fmt.Println("I'm ready")
	http.ListenAndServe(listenAddr, nil)
}
