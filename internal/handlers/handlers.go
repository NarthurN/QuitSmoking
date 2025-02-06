package handlers

import "net/http"


func Home(w http.ResponseWriter, r *http.Request) {
	msg := "Hello!"
	w.Write([]byte(msg))
}