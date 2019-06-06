package resthandler

import (
	"github.com/bitly/go-simplejson"
	"log"
	"net/http"
)

// makeErrorResponse generates an http response with the message an the code specified
func MakeErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := simplejson.New()
	response.Set("error", message)
	responsePayload, err := response.MarshalJSON()
	if err != nil {
		log.Panicln(err)
	}
	w.Write(responsePayload)
}
