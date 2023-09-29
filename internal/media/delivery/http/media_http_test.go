package http_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"source.toby3d.me/toby3d/pub/internal/common"
	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/media"
	delivery "source.toby3d.me/toby3d/pub/internal/media/delivery/http"
)

func TestHandler_Upload(t *testing.T) {
	t.Parallel()

	testConfig := domain.TestConfig(t)
	testFile := domain.TestFile(t)

	buf := bytes.NewBuffer(nil)
	form := multipart.NewWriter(buf)
	formWriter, err := form.CreateFormFile("file", "photo.jpg")
	if err != nil {
		t.Fatal(err)
	}

	if _, err = formWriter.Write(testFile.Content); err != nil {
		t.Fatal(err)
	}

	if err = form.Close(); err != nil {
		t.Fatal(err)
	}

	expect := testConfig.HTTP.BaseURL().JoinPath("media", "abc123"+testFile.Ext())

	req := httptest.NewRequest(http.MethodPost, "https://media.example.com", buf)
	req.Header.Set(common.HeaderContentType, form.FormDataContentType())

	w := httptest.NewRecorder()
	delivery.NewHandler(
		media.NewStubUseCase(nil, testFile, expect), *testConfig).
		ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("%s %s = %d, want %d", req.Method, req.RequestURI, resp.StatusCode, http.StatusCreated)
	}

	if location := resp.Header.Get(common.HeaderLocation); location != expect.String() {
		t.Errorf("%s %s = %s, want not empty", req.Method, req.RequestURI, location)
	}
}

func TestHandler_Download(t *testing.T) {
	t.Parallel()

	testConfig := domain.TestConfig(t)
	testFile := domain.TestFile(t)

	req := httptest.NewRequest(http.MethodGet, "https://media.example.com/"+testFile.LogicalName(), nil)
	w := httptest.NewRecorder()

	delivery.NewHandler(
		media.NewStubUseCase(nil, testFile, nil), *testConfig).
		ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("%s %s = %d, want %d", req.Method, req.RequestURI, resp.StatusCode, http.StatusOK)
	}

	contentType, mediaType := resp.Header.Get(common.HeaderContentType), testFile.MediaType()
	if contentType != mediaType {
		t.Errorf("%s %s = '%s', want '%s'", req.Method, req.RequestURI, contentType, mediaType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(body, testFile.Content) {
		t.Error("stored and received file contents is not the same")
	}
}
