package handlers

import (
	"encoding/json"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockLister struct {
	videos        []*domain.Video
	limit, offset uint32
}

func (l *mockLister) List(offset, limit uint32, handler func(v *domain.Video) (bool, error)) error {
	l.offset = offset
	l.limit = limit
	for _, v := range l.videos {
		if ok, err := handler(v); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

func TestList(t *testing.T) {
	w := httptest.NewRecorder()
	lister := mockLister{}
	lister.videos = []*domain.Video{
		domain.MakeVideo("video-id1", "video1-name", "video1-url", "", 13, domain.StatusUploaded),
		domain.MakeVideo("video-id2", "video1 name 2", "video2-path", "video2-thumb", 42, domain.StatusReady),
	}

	r := httptest.NewRequest("GET", "/api/v1/list?limit=3&skip=1", nil)
	list(&lister, w, r)

	if lister.limit != 3 {
		t.Errorf("Invalid limit. Have: %d, want: %d", lister.limit, 3)
	}
	if lister.offset != 1 {
		t.Errorf("Invalid offset. Have: %d, want: %d", lister.offset, 1)
	}

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
