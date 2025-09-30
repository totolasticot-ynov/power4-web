package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bienvenue sur mon API Go !")
		fmt.Fprintln(w, "HEHE!")

	})

	fmt.Println("url : http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
