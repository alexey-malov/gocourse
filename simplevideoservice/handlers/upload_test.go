package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"
)

type uploadRequest struct {
	formName    string
	fileName    string
	contentType string
}

func TestUploadGoodFile(t *testing.T) {
	runUploadTest(t, uploadRequest{fileName: "video.mp4", formName: "file[]", contentType: "video/mp4"}, http.StatusOK)
}

func TestUploadWrongContentType(t *testing.T) {
	runUploadTest(t, uploadRequest{fileName: "video.mp4", formName: "file[]", contentType: "image/png"}, http.StatusBadRequest)
}

func TestUploadWrongFormName(t *testing.T) {
	runUploadTest(t, uploadRequest{fileName: "video.mp4", formName: "thefile[]", contentType: "video/mp4"}, http.StatusBadRequest)
}
func runUploadTest(t *testing.T, ur uploadRequest, expectedStatus int) {
	u := mockUploader{}
	h := handlerBase{&u, nil}

	fileContent := "fileContent"
	r, err := makeUploadRequest(ur.fileName, fileContent, ur.contentType, ur.formName)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	h.upload(w, r)

	if w.Code != expectedStatus {
		t.Errorf("Invalid HTTP status. Want %d, got %d", expectedStatus, w.Code)
	}

	if w.Code != http.StatusOK {
		return
	}
	if u.content != fileContent {
		t.Errorf(`Wrong fileContent. Want: "%s", got:"%s"`, fileContent, u.content)
	}
	if u.fileName != ur.fileName {
		t.Errorf(`Wrong file fileName. Want: "%s", got: "%s"`, ur.fileName, u.fileName)
	}
}

func makeUploadRequest(fileName, content, contentType, formName string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s";`, formName, fileName))
	hdr.Set("Content-Type", contentType)
	part, err := writer.CreatePart(hdr)
	if err != nil {
		return nil, err
	}
	io.Copy(part, strings.NewReader(content))
	if err = writer.Close(); err != nil {
		return nil, err
	}

	r := httptest.NewRequest("POST", "/video", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	return r, nil
}
