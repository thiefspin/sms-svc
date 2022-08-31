package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type SMSRequest struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

func (sms *SMSRequest) valid() bool {
	return sms.Sender != "" && sms.Receiver != ""
}

var smsQueue = make(chan SMSRequest, math.MaxInt16)

func respond(o any, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	withJsonEncoding(o, func(b *bytes.Buffer) {
		fmt.Fprintf(w, b.String())
	})
}

func withJsonDecoding(body io.ReadCloser, target any, w http.ResponseWriter, f func()) {
	jsonErr := json.NewDecoder(body).Decode(&target)
	if jsonErr != nil {
		res := APIError{StatusCode: 400, Message: "Invalid JSON input"}
		respond(res, 400, w)
	} else {
		f()
	}
}

func withJsonEncoding(o any, f func(b *bytes.Buffer)) {
	jsonObject, err := json.Marshal(o)
	if err != nil {
		log.Println(err)
		return
	}
	f(bytes.NewBuffer(jsonObject))
}

func callSMSGateway(sms SMSRequest) bool {
	success := true
	withJsonEncoding(sms, func(b *bytes.Buffer) {
		resp, err := http.Post("http://localhost:4000", "application/json", b)
		if err != nil || resp.StatusCode == 429 {
			log.Println("Failed to call SMS API")
			success = false
			if err != nil {
				log.Println(err)
			}
		}
	})
	return success
}

func create(w http.ResponseWriter, r *http.Request) {
	var sms SMSRequest
	withJsonDecoding(r.Body, &sms, w, func() {
		if !sms.valid() {
			res := APIError{StatusCode: 400, Message: "Invalid JSON input"}
			respond(res, 400, w)
		} else {
			log.Println("Server received SMS request: ", sms)
			smsQueue <- sms
			respond(sms, 202, w)
		}
	})
}

func handleMessages() {
	queueSize := len(smsQueue)
	if queueSize > 0 {
		sms := <-smsQueue
		if !callSMSGateway(sms) {
			log.Print("Failed. Waiting before trying again....")
			smsQueue <- sms
			time.Sleep(60 * time.Second) //No sense in retrying immediately
		}
		log.Println("SMS's in the queue: ", queueSize)
	}

}

func scheduler() {
	for {
		handleMessages()
	}
}

func main() {
	go func() {
		scheduler()
	}()
	r := mux.NewRouter()
	r.HandleFunc("/api/sms", create).Methods("POST")
	log.Fatal(http.ListenAndServe(":8082", r))

}
