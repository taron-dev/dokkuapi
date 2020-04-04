package helper

import (
	"encoding/json"
	"net/http"
)

func RespondWithData(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			w.Write([]byte(""))
		}
	}
}

func RespondWithMessage(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	jsonData, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		w.Write([]byte(""))
	}
	w.Write(jsonData)

}
