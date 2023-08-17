package user

import "net/http"

// Dashboard is a simple handler for user app
func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is user dashboard"))
}
