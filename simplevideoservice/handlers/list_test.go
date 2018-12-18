package handlers

import (
	"encoding/json"
	"github.com/alexey-malov/gocourse/simplevideoservice/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestList(t *testing.T) {
	w := httptest.NewRecorder()
	vr := &mockRepo{}
	vr.videos = []model.VideoItem{
		model.MakeVideoItem("video-id1", "video1-name", 13),
		model.MakeVideoItem("video-id2", "video1 name 2", 42),
	}
	list(vr, w, nil)
	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have: %d, want: %d.", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	items := make([]videoListItem, 10)
	if err = json.Unmarshal(jsonString, &items); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}

	if len(items) != 2 {
		t.Error("3 list items expected")
	}
}
