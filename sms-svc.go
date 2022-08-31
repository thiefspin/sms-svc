package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"sms-svc/restutils"
	"time"
)

type SMSRequest struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

func (sms *SMSRequest) valid() bool {
	return sms.Sender != "" && sms.Receiver != ""
}

var smsQueue = make(chan SMSRequest, math.MaxInt16)

func callSMSGateway(sms SMSRequest) bool {
	success := true
	utils.WithJsonEncoding(sms, func(b *bytes.Buffer) {
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
	utils.WithJsonDecoding(r.Body, &sms, w, func() {
		if !sms.valid() {
			res := utils.APIError{StatusCode: 400, Message: "Invalid JSON input"}
			utils.Respond(res, 400, w)
		} else {
			log.Println("Server received SMS request: ", sms)
			smsQueue <- sms
			utils.Respond(sms, 202, w)
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
