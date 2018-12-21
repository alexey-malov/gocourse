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

type mockRepo struct {
	videos     []domain.Video
	findResult *domain.Video
}

func (r *mockRepo) Enumerate(handler func(v domain.Video) bool) error {
	for _, v := range r.videos {
		if !handler(v) {
			return nil
		}
	}
	return nil
}

func (r *mockRepo) Find(id string) (*domain.Video, error) {
	return r.findResult, nil
}

func (r *mockRepo) Add(v domain.Video) error {
	r.videos = append(r.videos, v)
	return nil
}

func TestVideo(t *testing.T) {
	w := httptest.NewRecorder()
	v := domain.MakeVideo("video-id", "video-name", 12345)
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
