package http_test

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"

	"source.toby3d.me/toby3d/pub/internal/entry/delivery/http"
)

type testRequest struct {
	Delete  *http.Delete   `json:"delete,omitempty"`
	Content []http.Content `json:"content,omitempty"`
	Photo   []*http.Figure `json:"photo,omitempty"`
}

func TestRequest(t *testing.T) {
	t.Parallel()

	req := new(http.Request)
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
		out  http.Content
	}{{
		name: "plain",
		in:   `"Hello World"`,
		out: http.Content{
			HTML:  nil,
			Value: "Hello World",
		},
	}, {
		name: "html",
		in:   `{"html":"<b>Hello</b> <i>World</i>"}`,
		out: http.Content{
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

			if out == nil || len(out.Content) == 0 {
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
		in   http.Content
		out  string
		name string
	}{{
		name: "plain",
		in: http.Content{
			HTML:  nil,
			Value: `Hello World`,
		},
		out: `{"content":["Hello World"]}`,
	}, {
		name: "html",
		in: http.Content{
			HTML:  testContent,
			Value: "",
		},
		out: `{"content":[{"html":"\u003cb\u003eHello\u003c/b\u003e \u003ci\u003eWorld\u003c/i\u003e"}]}`,
	}, {
		name: "both",
		in: http.Content{
			HTML:  testContent,
			Value: `Hello World`,
		},
		out: `{"content":[{"html":"\u003cb\u003eHello\u003c/b\u003e \u003ci\u003eWorld\u003c/i\u003e"}]}`,
	}} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := json.Marshal(testRequest{
				Content: []http.Content{tc.in},
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
		out  http.Delete
	}{{
		name: "values",
		in:   `{"category":["indieweb"]}`,
		out: http.Delete{
			Keys: nil,
			Values: http.Properties{
				Category: []string{"indieweb"},
			},
		},
	}, {
		name: "keys",
		in:   `["category"]`,
		out: http.Delete{
			Keys:   []string{"category"},
			Values: http.Properties{},
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
		out  http.Figure
	}{{
		name: "alt",
		in:   `{"value":"https://photos.example.com/globe.gif","alt":"Spinning globe animation"}`,
		out: http.Figure{
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
		out: http.Figure{
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
