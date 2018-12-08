package handlers

import (
	"fmt"
	"net/http"
)

func list(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, `[{
    "id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "name": "Black Retrospetive Woman",
    "duration": 15,
    "thumbnail": "/content/d290f1ee-6c54-4b01-90e6-d701748f0851/screen.jpg"
}]`)
	if err == nil {
		return
	}
}
