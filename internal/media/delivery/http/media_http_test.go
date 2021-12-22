package http_test

import (
	"bytes"
	"embed"
	"mime/multipart"
	"path/filepath"
	"sync"
	"testing"

	"github.com/fasthttp/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	http "github.com/valyala/fasthttp"

	"source.toby3d.me/website/micropub/internal/domain"
	delivery "source.toby3d.me/website/micropub/internal/media/delivery/http"
	repository "source.toby3d.me/website/micropub/internal/media/repository/memory"
	"source.toby3d.me/website/micropub/internal/media/usecase"
	"source.toby3d.me/website/micropub/internal/testing/httptest"
)

type TestCase struct {
	name           string
	fileName       string
	expContentType string
}

//go:embed testdata/*
var testData embed.FS

func TestUpload(t *testing.T) {
	t.Parallel()

	cfg := domain.TestConfig(t)
	r := router.New()
	delivery.New(cfg, usecase.NewMediaUseCase(repository.NewMemoryMediaRepository(new(sync.Map)))).Register(r)

	client, _, cleanup := httptest.New(t, r.Handler)
	t.Cleanup(cleanup)

	for _, testCase := range []TestCase{
		{
			name:           "jpg",
			fileName:       "sunset.jpg",
			expContentType: "image/jpeg",
		}, {
			name:           "png",
			fileName:       "micropub-rocks.png",
			expContentType: "image/png",
		}, {
			name:           "gif",
			fileName:       "w3c-socialwg.gif",
			expContentType: "image/gif",
		},
	} {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			upResp := upload(t, client, testCase.fileName)
			defer http.ReleaseResponse(upResp)

			assert.Equal(t, upResp.StatusCode(), http.StatusCreated, "returned HTTP 201")
			assert.NotNil(t, upResp.Header.Peek(http.HeaderLocation), "returned a Location header")

			downResp := download(t, client, upResp.Header.Peek(http.HeaderLocation))
			assert.Equal(t, http.StatusOK, downResp.StatusCode(), "the URL exists")
			assert.Equal(t, testCase.expContentType, string(downResp.Header.ContentType()),
				"has the expected content type")
		})
	}
}

func upload(tb testing.TB, client *http.Client, fileName string) *http.Response {
	tb.Helper()

	contents, err := testData.ReadFile(filepath.Join("testdata", fileName))
	require.NoError(tb, err)

	// NOTE(toby3d): upload
	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)

	ff, err := w.CreateFormFile("file", fileName)
	require.NoError(tb, err)

	_, err = ff.Write(contents)
	require.NoError(tb, err)
	require.NoError(tb, w.Close())

	req := http.AcquireRequest()
	defer http.ReleaseRequest(req)
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetMultipartFormBoundary(w.Boundary())
	req.SetRequestURI("https://example.com/media")
	req.SetBody(buf.Bytes())

	resp := http.AcquireResponse()
	require.NoError(tb, client.Do(req, resp))

	return resp
}

func download(tb testing.TB, client *http.Client, location []byte) *http.Response {
	tb.Helper()

	req := http.AcquireRequest()
	defer http.ReleaseRequest(req)
	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURIBytes(location)

	resp := http.AcquireResponse()
	require.NoError(tb, client.Do(req, resp))

	return resp
}
