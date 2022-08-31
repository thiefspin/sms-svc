# SMS Service
Provides a HTTP REST interface to a SMS gateway.
* Developed on Golang 1.19

## Running the application

### Running the mock gateway
Before running the mock gateway install this package via `go get "golang.org/x/time/rate"`. 

Afterwards just build the application with `go build sms-gateway.go` and run with `./sms-gateway`

### Running the SMS Service
Get the `mux` dependency with `"github.com/gorilla/mux"` and simply build the application with `go build sms-svc.go` and run with `./sms-svc`
Create a SMS in the following way:
```shell
curl -v -X POST http://localhost:8082/api/sms -H 'Content-Type: application/json' -d
 '{"sender": "User1", "receiver": "User2", "message": "This is sms number 1"}'
```

## Extra Information
### Possible future enhancements
* SMS messages are not guaranteed to be sent to the gateway in order in the event that there are retries to the gateway
* Might be a better way to handle the async sending of messages
