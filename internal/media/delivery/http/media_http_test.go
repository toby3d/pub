package http_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"mime/multipart"
	"sync"
	"testing"

	"github.com/fasthttp/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	http "github.com/valyala/fasthttp"

	delivery "source.toby3d.me/website/micropub/internal/media/delivery/http"
	repository "source.toby3d.me/website/micropub/internal/media/repository/memory"
	"source.toby3d.me/website/micropub/internal/media/usecase"
	"source.toby3d.me/website/micropub/internal/testing/httptest"
)

const testFileName string = "sunset.jpg"

func TestUpload(t *testing.T) {
	t.Parallel()

	_, contents := testFile(t)
	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)

	ff, err := w.CreateFormFile("file", testFileName)
	require.NoError(t, err)

	_, err = ff.Write(contents)
	require.NoError(t, err)

	require.NoError(t, w.Close())

	r := router.New()
	r.POST("/media", delivery.New(usecase.NewMediaUseCase(repository.NewMemoryMediaRepository(new(sync.Map)))).
		Update)

	client, _, cleanup := httptest.New(t, r.Handler)
	t.Cleanup(cleanup)

	req := http.AcquireRequest()
	defer http.ReleaseRequest(req)
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetMultipartFormBoundary(w.Boundary())
	req.SetRequestURI("https://example.com/media")
	req.SetBody(buf.Bytes())

	resp := http.AcquireResponse()
	defer http.ReleaseResponse(resp)

	require.NoError(t, client.Do(req, resp))
	assert.Equal(t, resp.StatusCode(), http.StatusCreated)
	require.NotNil(t, resp.Header.Peek(http.HeaderLocation))
}

func TestDownload(t *testing.T) {
	t.Parallel()

	fileName, contents := testFile(t)

	repo := repository.NewMemoryMediaRepository(new(sync.Map))
	require.NoError(t, repo.Create(context.Background(), fileName, contents))

	r := router.New()
	r.GET("/media/{fileName:*}", delivery.New(usecase.NewMediaUseCase(repo)).Read)

	client, _, cleanup := httptest.New(t, r.Handler)
	t.Cleanup(cleanup)

	status, body, err := client.Get(nil, "https://example.com/media/"+fileName)
	assert.NoError(t, err)
	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, contents, body)
}

func testFile(tb testing.TB) (string, []byte) {
	tb.Helper()

	fileName := make([]byte, usecase.DefaultNameLength)
	_, err := rand.Read(fileName)
	require.NoError(tb, err)

	contents := make([]byte, 128)
	_, err = rand.Read(contents)
	require.NoError(tb, err)

	return base64.RawURLEncoding.EncodeToString(fileName) + ".jpg", contents
}
