package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockLister struct {
	videos        []*domain.Video
	limit, offset uint32
}

func (l *mockLister) List(offset, limit uint32, handler func(v *domain.Video) (bool, error)) error {
	if l.videos == nil {
		return errors.New("Simulated errror")
	}
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

type listResponse struct {
	status        int
	offset, limit uint32
}

func (l *mockLister) expectedJSON() []videoListItem {
	result := make([]videoListItem, 0)
	for _, v := range l.videos {
		result = append(result, makeVideoListItem(*v))
	}
	return result
}

func testListImpl(t *testing.T, videos []*domain.Video, query string, expectedResponse listResponse) {
	w := httptest.NewRecorder()
	lister := &mockLister{}
	lister.videos = videos

	r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/list?%s", query), nil)
	list(lister, w, r)

	response := w.Result()
	if response.StatusCode != expectedResponse.status {
		t.Errorf("Status code is wrong. Have: %d, want: %d.", response.StatusCode, expectedResponse.status)
	}
	if response.StatusCode != http.StatusOK {
		return
	}

	if lister.limit != expectedResponse.limit {
		t.Errorf("Invalid limit. Have: %d, want: %d", lister.limit, expectedResponse.limit)
	}
	if lister.offset != expectedResponse.offset {
		t.Errorf("Invalid offset. Have: %d, want: %d", lister.offset, expectedResponse.offset)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err = response.Body.Close(); err != nil {
		t.Fatal(err)
	}

	jsonItems := make([]videoListItem, 0)
	if err = json.Unmarshal(jsonString, &jsonItems); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}

	if !reflect.DeepEqual(jsonItems, lister.expectedJSON()) {
		t.Errorf("Invalid JSON items")
	}
}

func TestList(t *testing.T) {
	testListImpl(t, []*domain.Video{
		domain.MakeVideo("video-id1", "video1-name", "video1-url", "", 13, domain.StatusUploaded),
		domain.MakeVideo("video-id2", "video1 name 2", "video2-path", "video2-thumb", 42, domain.StatusReady),
	}, "limit=3&skip=1", listResponse{http.StatusOK, 1, 3})
}

func TestListWithoutOffset(t *testing.T) {
	testListImpl(t, nil, "limit=3",
		listResponse{http.StatusBadRequest, 0, 0})
}

func TestListWithoutLimit(t *testing.T) {
	testListImpl(t, nil, "skip=2",
		listResponse{http.StatusBadRequest, 0, 0})
}

func TestListError(t *testing.T) {
	testListImpl(t, nil, "skip=2&limit=3",
		listResponse{http.StatusInternalServerError, 0, 0})
}
