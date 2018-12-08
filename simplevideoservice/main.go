package main

import (
	"fmt"
	"github.com/alexey-malov/gocourse/simplevideoservice/handlers"
	"net/http"
)

func main() {
	router := handlers.Router()
	fmt.Println(http.ListenAndServe(":8000", router))
}
