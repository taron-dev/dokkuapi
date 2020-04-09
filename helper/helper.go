package helper

import (
	"encoding/json"
	"net/http"
)

// Response serves to unify inner method return value
type Response struct {
	Value   interface{}
	Status  int
	Message string
}

// RespondWithData unify response format with data included
func RespondWithData(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			w.Write([]byte(""))
		}
	}
}

// RespondWithMessage unify response format with message instead of data
func RespondWithMessage(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	jsonData, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		w.Write([]byte(""))
	}
	w.Write(jsonData)

}

// Decode body from request into object
func Decode(w http.ResponseWriter, r *http.Request, object interface{}) error {
	return json.NewDecoder(r.Body).Decode(object)
}
