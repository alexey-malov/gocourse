package handlers

import (
	"encoding/json"
	"net/http"
)

type videoContent struct {
	videoJson
	Url string `json:"url"`
}

func video(w http.ResponseWriter, _ *http.Request) {
	v := videoContent{
		videoJson{
			"d290f1ee-6c54-4b01-90e6-d701748f0851",
			"Black Retrospetive Woman",
			15,
			"/content/d290f1ee-6c54-4b01-90e6-d701748f0851/screen.jpg",
		},
		"/content/d290f1ee-6c54-4b01-90e6-d701748f0851/index.mp4",
	}

	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		return
	}

	if err == nil {
		return
	}
}
