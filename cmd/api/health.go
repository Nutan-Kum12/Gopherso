package main

import "net/http"

func (app application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
	// app.store.Users.Create(r.Context()) //<- Example usage of the store
	// app.store.Posts.Create(r.Context()) //<- Example usage of the store
}
