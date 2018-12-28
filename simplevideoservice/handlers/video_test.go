package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockFinder struct {
	video *domain.Video
	error error
}

func (f *mockFinder) Find(id string) (*domain.Video, error) {
	return f.video, f.error
}

type videoResponse struct {
	status int
	body   *videoContent
}

func TestSuccessfulVideoInfoRequest(t *testing.T) {
	v := domain.MakeVideo("id", "name", "url", "thumbnail", 42, domain.StatusProcessing)
	expectedResponse := videoResponse{
		http.StatusOK,
		&videoContent{
			videoListItem{v.Id(), v.Name(), v.Duration(), v.ThumbnailURL()},
			v.VideoUrl()}}

	testVideoImpl(t, expectedResponse.body.Id, mockFinder{v, nil}, expectedResponse)
}

func TestMissingVideoInfo(t *testing.T) {
	testVideoImpl(t, "missing-id", mockFinder{nil, nil}, videoResponse{http.StatusNotFound, nil})
}

func TestVideoInfoWithoutId(t *testing.T) {
	testVideoImpl(t, "", mockFinder{nil, nil}, videoResponse{http.StatusBadRequest, nil})
}

func TestVideoInfoSearchingError(t *testing.T) {
	testVideoImpl(t, "id",
		mockFinder{nil, errors.New("Simulated error during finder.Find")},
		videoResponse{http.StatusInternalServerError, nil})
}

func testVideoImpl(t *testing.T, id string, finder mockFinder, response videoResponse) {
	w := httptest.NewRecorder()
	h := &UseCases{finder: &finder}
	r := httptest.NewRequest("GET", fmt.Sprintf("/video/%s", id), nil)
	if id != "" {
		vars := map[string]string{"ID": id}
		r = mux.SetURLVars(r, vars)
	}

	h.video(w, r)

	res := w.Result()
	if res.StatusCode != response.status {
		t.Errorf("Status code is wrong. Got: %d, want: %d.", res.StatusCode, response.status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err = res.Body.Close(); err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		return
	}

	desc := videoContent{}
	if err = json.Unmarshal(body, &desc); err != nil {
		t.Errorf(`Failed to unmarshal videoContent from "%s""`, string(body))
	}

	if response.body == nil {
		t.Fatal("Invalid response")
	}

	if desc != *response.body {
		t.Errorf("Invalid JSON content. Got %v, want %v", desc, *response.body)
	}
}
