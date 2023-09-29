package http_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"

	"source.toby3d.me/toby3d/pub/internal/common"
	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/entry"
	delivery "source.toby3d.me/toby3d/pub/internal/entry/delivery/http"
	"source.toby3d.me/toby3d/pub/internal/media"
)

type testRequest struct {
	Delete  *delivery.Delete   `json:"delete,omitempty"`
	Content []delivery.Content `json:"content,omitempty"`
	Photo   []*delivery.Figure `json:"photo,omitempty"`
}

func TestHandler_Create(t *testing.T) {
	t.Parallel()

	t.Run("form", func(t *testing.T) {
		t.Parallel()

		for name, input := range map[string]url.Values{
			"simple": {
				"h":       []string{"entry"},
				"content": []string{"Micropub test of creating a basic h-entry"},
			},
			"categories": {
				"h": []string{"entry"},
				"content": []string{"Micropub test of creating an h-entry with categories. " +
					"This post should have two categories, test1 and test2"},
				"category[]": []string{"test1", "test2"},
			},
			"category": {
				"h": []string{"entry"},
				"content": []string{"Micropub test of creating an h-entry with one category. " +
					"This post should have one category, test1"},
				"category": []string{"test1"},
			},
		} {
			name, input := name, input

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				doCreateRequest(t, strings.NewReader(input.Encode()),
					common.MIMEApplicationFormCharsetUTF8)
			})
		}
	})

	t.Run("json", func(t *testing.T) {
		t.Parallel()

		for name, input := range map[string]string{
			"simple":     `{"type": ["h-entry"], "properties": {"content": ["Micropub test of creating an h-entry with a JSON request"]}}`,
			"categories": `{"type": ["h-entry"], "properties": {"content": ["Micropub test of creating an h-entry with a JSON request containing multiple categories. This post should have two categories, test1 and test2."], "category": ["test1", "test2"]}}`,
			"html":       `{"type": ["h-entry"], "properties": {"content": [{"html": "<p>This post has <b>bold</b> and <i>italic</i> text.</p>"}]}}`,
			"photo":      `{"type": ["h-entry"], "properties": {"content": ["Micropub test of creating a photo referenced by URL. This post should include a photo of a sunset."], "photo": ["https://micropub.rocks/media/sunset.jpg"]}}`,
			"object":     `{"type": ["h-entry"], "properties": {"published": ["2017-05-31T12:03:36-07:00"], "content": ["Lunch meeting"], "checkin": [{"type": ["h-card"], "properties": {"name": ["Los Gorditos"], "url": ["https://foursquare.com/v/502c4bbde4b06e61e06d1ebf"], "latitude": [45.524330801154], "longitude": [-122.68068808051], "street-address": ["922 NW Davis St"], "locality": ["Portland"], "region": ["OR"], "country-name": ["United States"], "postal-code": ["97209"]}}]}}`,
			"photo-alt":  `{"type": ["h-entry"], "properties": {"content": ["Micropub test of creating a photo referenced by URL with alt text. This post should include a photo of a sunset."], "photo": [{"value": "https://micropub.rocks/media/sunset.jpg", "alt": "Photo of a sunset"}]}}`,
			"photos":     `{"type": ["h-entry"], "properties": {"content": ["Micropub test of creating multiple photos referenced by URL. This post should include a photo of a city at night."], "photo": ["https://micropub.rocks/media/sunset.jpg", "https://micropub.rocks/media/city-at-night.jpg"]}}`,
		} {
			name, input := name, input

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				doCreateRequest(t, strings.NewReader(input), common.MIMEApplicationJSONCharsetUTF8)
			})
		}
	})

	// TODO(toby3d): multipart requests
}

func doCreateRequest(tb testing.TB, r io.Reader, contentType string) {
	tb.Helper()

	req := httptest.NewRequest(http.MethodPost, "https://example.com/", r)
	req.Header.Set(common.HeaderContentType, contentType)

	w := httptest.NewRecorder()
	delivery.NewHandler(entry.NewStubUseCase(nil, domain.TestEntry(tb), true),
		media.NewDummyUseCase()).ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		tb.Errorf("%s %s = %d, expect %d or %d", req.Method, req.RequestURI, resp.StatusCode,
			http.StatusCreated, http.StatusAccepted)
	}

	if location := resp.Header.Get(common.HeaderLocation); location == "" {
		tb.Errorf("%s %s = returns empty Location header, want non-empty", req.Method, req.RequestURI)
	}
}

func TestRequest(t *testing.T) {
	t.Parallel()

	req := new(delivery.Request)
	if err := json.NewDecoder(strings.NewReader(`{
		  "action": "update",
		  "url": "http://example.com/",
		  "add": {
		    "syndication": ["http://web.archive.org/web/20040104110725/https://aaronpk.example/2014/06/01/9/indieweb"]
		  }
		}`)).Decode(req); err != nil {
		t.Fatal(err)
	}

	if req.Action != "update" {
		t.Errorf("got %s, want %s", req.Action, "update")
	}
}

func TestContent_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	testContent, err := html.Parse(strings.NewReader(`<b>Hello</b> <i>World</i>`))
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		in   string
		out  delivery.Content
	}{{
		name: "plain",
		in:   `"Hello World"`,
		out: delivery.Content{
			HTML:  nil,
			Value: "Hello World",
		},
	}, {
		name: "html",
		in:   `{"html":"<b>Hello</b> <i>World</i>"}`,
		out: delivery.Content{
			HTML:  testContent,
			Value: "",
		},
	}} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := new(testRequest)
			if err := json.Unmarshal([]byte(`{"content": [`+tc.in+`]}`), out); err != nil {
				t.Fatal(err)
			}

			if len(out.Content) == 0 {
				t.Error("empty content result, want not nil")

				return
			}

			if diff := cmp.Diff(out.Content[0], tc.out); diff != "" {
				t.Errorf("%+s", diff)
			}
		})
	}
}

func TestContent_MarshalJSON(t *testing.T) {
	t.Parallel()

	testContent, err := html.Parse(strings.NewReader(`<b>Hello</b> <i>World</i>`))
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		in   delivery.Content
		out  string
		name string
	}{{
		name: "plain",
		in: delivery.Content{
			HTML:  nil,
			Value: `Hello World`,
		},
		out: `{"content":["Hello World"]}`,
	}, {
		name: "html",
		in: delivery.Content{
			HTML:  testContent,
			Value: "",
		},
		out: `{"content":[{"html":"\u003cb\u003eHello\u003c/b\u003e \u003ci\u003eWorld\u003c/i\u003e"}]}`,
	}, {
		name: "both",
		in: delivery.Content{
			HTML:  testContent,
			Value: `Hello World`,
		},
		out: `{"content":[{"html":"\u003cb\u003eHello\u003c/b\u003e \u003ci\u003eWorld\u003c/i\u003e"}]}`,
	}} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := json.Marshal(testRequest{
				Content: []delivery.Content{tc.in},
			})
			if err != nil {
				t.Fatal(err)
			}

			if string(out) != tc.out {
				t.Errorf("got '%s', want '%s'", out, tc.out)
			}
		})
	}
}

func TestDelete_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		in   string
		out  delivery.Delete
	}{{
		name: "values",
		in:   `{"category":["indieweb"]}`,
		out: delivery.Delete{
			Keys: nil,
			Values: delivery.Properties{
				Category: []string{"indieweb"},
			},
		},
	}, {
		name: "keys",
		in:   `["category"]`,
		out: delivery.Delete{
			Keys:   []string{"category"},
			Values: delivery.Properties{},
		},
	}} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := new(testRequest)
			if err := json.Unmarshal([]byte(`{"delete":`+tc.in+`}`), out); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(*out.Delete, tc.out); diff != "" {
				t.Errorf("%+s", diff)
			}
		})
	}
}

func TestFigure_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		in   string
		out  delivery.Figure
	}{{
		name: "alt",
		in:   `{"value":"https://photos.example.com/globe.gif","alt":"Spinning globe animation"}`,
		out: delivery.Figure{
			Alt: "Spinning globe animation",
			Value: &url.URL{
				Scheme: "https",
				Host:   "photos.example.com",
				Path:   "/globe.gif",
			},
		},
	}, {
		name: "plain",
		in:   `"https://photos.example.com/592829482876343254.jpg"`,
		out: delivery.Figure{
			Alt: "",
			Value: &url.URL{
				Scheme: "https",
				Host:   "photos.example.com",
				Path:   "/592829482876343254.jpg",
			},
		},
	}} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := new(testRequest)
			if err := json.Unmarshal([]byte(`{"photo":[`+tc.in+`]}`), out); err != nil {
				t.Fatal(err)
			}

			if len(out.Photo) == 0 {
				t.Fatal("empty photo value, want not nil")
			}

			if diff := cmp.Diff(out.Photo[0], &tc.out); diff != "" {
				t.Errorf("%+s", diff)
			}
		})
	}
}
