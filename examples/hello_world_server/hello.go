package main

import "net/http"

func main() {
	http.HandleFunc("/", sayHello)
	http.ListenAndServe(":8080", nil)
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!\n"))
}
