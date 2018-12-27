package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVideo(t *testing.T) {
	w := httptest.NewRecorder()
	v := domain.MakeVideo("video-id", "video-name", "video-url", "video-thumbnail", 12345, domain.StatusReady)
	r := httptest.NewRequest("GET", fmt.Sprintf("/video/%s", v.Id()), nil)
	vars := map[string]string{"ID": v.Id()}
	r = mux.SetURLVars(r, vars)

	videos := mockVideos{}
	h := &handlerBase{nil, &videos}

	videos.findResult = v
	h.video(w, r)

	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Got: %d, want: %d.", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err = response.Body.Close(); err != nil {
		t.Fatal(err)
	}

	videoDesc := videoContent{}
	if err = json.Unmarshal(jsonString, &videoDesc); err != nil {
		t.Errorf("Failed to unmarshal videoContent from %s", string(jsonString))
	}

	expectedDesc := makeVideoContent(*v)
	if expectedDesc != videoDesc {
		t.Errorf("Invalid JSON content. Got %v, want: %v", videoDesc, expectedDesc)
	}
}
