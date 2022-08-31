package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func WithJsonDecoding(body io.ReadCloser, target any, w http.ResponseWriter, f func()) {
	jsonErr := json.NewDecoder(body).Decode(&target)
	if jsonErr != nil {
		res := APIError{StatusCode: 400, Message: "Invalid JSON input"}
		Respond(res, 400, w)
	} else {
		f()
	}
}

func WithJsonEncoding(o any, f func(b *bytes.Buffer)) {
	jsonObject, err := json.Marshal(o)
	if err != nil {
		log.Println(err)
		return
	}
	f(bytes.NewBuffer(jsonObject))
}

func Respond(o any, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	WithJsonEncoding(o, func(b *bytes.Buffer) {
		fmt.Fprintf(w, b.String())
	})
}
