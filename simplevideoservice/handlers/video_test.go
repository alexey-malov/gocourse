package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/alexey-malov/gocourse/simplevideoservice/model"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRepo struct {
	videos     []model.VideoItem
	findResult *model.VideoItem
}

/*
var videoItems = []model.VideoItem{
	{"d290f1ee-6c54-4b01-90e6-d701748f0851", "Black Retrospetive Woman", 15},
	{"sldjfl34-dfgj-523k-jk34-5jk3j45klj34", "Go Rally TEASER-HD", 41},
	{"hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345", "Танцор", 92},
}

type videoItemCallback func(v videoItem) bool

func enumVideos(cb videoItemCallback) {
	for _, v := range videoItems {
		if !cb(v) {
			return
		}
	}
}*/

func (r *mockRepo) EnumVideos(handler func(v model.VideoItem) bool) error {
	for _, v := range r.videos {
		if !handler(v) {
			return nil
		}
	}
	return nil
}

func (r *mockRepo) FindVideo(id string) (*model.VideoItem, error) {
	return r.findResult, nil
}

func (r *mockRepo) AddVideo(v model.VideoItem) error {
	r.videos = append(r.videos, v)
	return nil
}

func TestVideo(t *testing.T) {
	w := httptest.NewRecorder()
	v := model.MakeVideoItem("video-id", "video-name", 12345)
	r := httptest.NewRequest("GET", fmt.Sprintf("/video/%s", v.Id()), nil)
	vars := map[string]string{"ID": v.Id()}
	r = mux.SetURLVars(r, vars)

	vr := &mockRepo{}
	vr.findResult = &v
	video(vr, w, r)

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

	expectedDesc := makeVideoContent(v)
	if expectedDesc != videoDesc {
		t.Errorf("Invalid JSON content. Got %v, want: %v", videoDesc, expectedDesc)
	}
}
