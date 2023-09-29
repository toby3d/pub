package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"source.toby3d.me/toby3d/pub/internal/common"
	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/entry"
	"source.toby3d.me/toby3d/pub/internal/media"
)

type (
	Handler struct {
		entries entry.UseCase
		media   media.UseCase
	}

	Request struct {
		Action string `json:"action"`
	}

	RequestCreate struct {
		Properties Properties `json:"properties"`
		Type       []string   `json:"type"` // h-entry
	}

	RequestSource struct {
		URL        URL
		Q          string
		Properties []string
	}

	RequestUpdate struct {
		Replace *Properties `json:"replace,omitempty"`
		Add     *Properties `json:"add,omitempty"`
		Delete  *Delete     `json:"delete,omitempty"`
		URL     URL         `json:"url"`
		Action  string      `json:"action"` // update
	}

	RequestDelete struct {
		URL    URL    `json:"url"`
		Action string `json:"action"`
	}

	RequestUndelete struct {
		URL    URL    `json:"url"`
		Action string `json:"action"`
	}

	ResponseSource struct {
		Properties Properties `json:"properties"`
		Type       []string   `json:"type,omitempty"`
	}

	Properties struct {
		Audio       []Figure   `json:"audio,omitempty"`
		Featured    []URL      `json:"featured,omitempty"`
		InReplyTo   []URL      `json:"in-reply-to,omitempty"`
		Like        []URL      `json:"like,omitempty"`
		Repost      []URL      `json:"repost,omitempty"`
		Syndication []URL      `json:"syndication,omitempty"`
		URL         []URL      `json:"url,omitempty"`
		Video       []Figure   `json:"video,omitempty"`
		Published   []DateTime `json:"published,omitempty"`
		Updated     []DateTime `json:"updated,omitempty"`
		Content     []Content  `json:"content,omitempty"`
		Photo       []Figure   `json:"photo,omitempty"`
		Category    []string   `json:"category,omitempty"`
		Name        []string   `json:"name,omitempty"`
		RSVP        []string   `json:"rsvp,omitempty"`
		Summary     []string   `json:"summary,omitempty"`
		UID         []string   `json:"uid,omitempty"`
		Altitude    []float32  `json:"altitude,omitempty"`
		Latitude    []float32  `json:"latitude,omitempty"`
		Longitude   []float32  `json:"longitude,omitempty"`
		Duration    []uint64   `json:"duration,omitempty"`
		Size        []uint64   `json:"size,omitempty"`
		// Author        []Author        `json:"author,omitempty"`
		// Location      []Location      `json:"location,omitempty"`
		// LikeOf        []LikeOf        `json:"like-of,omitempty"`
		// RepostOf      []RepostOf      `json:"repost-of,omitempty"`
		// Comment       []Comment       `json:"comment,omitempty"`
		// BookmarkOf    []BookmarkOf    `json:"bookmark-of,omitempty"`
		// ListenOf      []ListenOf      `json:"listen-of,omitempty"`
		// WatchOf       []WatchOf       `json:"watch-of,omitempty"`
		// ReadOf        []ReadOf        `json:"read-of,omitempty"`
		// TranslationOf []TranslationOf `json:"translation-of,omitempty"`
		// Checkin       []Checkin       `json:"checkin,omitempty"`
		// PlayOf        []PlayOf        `json:"play-of,omitempty"`
	}

	Delete struct {
		Keys   []string   `json:"-"`
		Values Properties `json:"-"`
	}

	Figure struct {
		Value *url.URL `json:"-"`
		Alt   string   `json:"-"`
	}

	Content struct {
		HTML  *html.Node `json:"-"`
		Value string     `json:"-"`
	}

	URL struct {
		*url.URL `json:"-"`
	}

	DateTime struct {
		time.Time `json:"-"`
	}

	Action struct {
		domain.Action `json:"-"`
	}

	bufferHTML struct {
		HTML string `json:"html,omitempty"`
	}

	bufferMedia struct {
		Value string `json:"value,omitempty"`
		Alt   string `json:"alt,omitempty"`
	}
)

const MaxBodySize int64 = 100 * 1024 * 1024 // 100mb

func NewHandler(entries entry.UseCase, media media.UseCase) *Handler {
	return &Handler{
		entries: entries,
		media:   media,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	case "", http.MethodGet:
		h.handleSource(w, r)
	case http.MethodPost:
		mediaType, _, err := mime.ParseMediaType(r.Header.Get(common.HeaderContentType))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		switch mediaType {
		default:
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		case common.MIMEApplicationJSON:
			buf := bytes.NewBuffer(nil)
			if _, err := buf.ReadFrom(r.Body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			req := new(Request)
			_ = json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(req)
			r.Body = io.NopCloser(buf)

			switch req.Action {
			default:
				h.handleCreate(w, r)
			case domain.ActionUpdate.String():
				h.handleUpdate(w, r)
			case domain.ActionDelete.String():
				h.handleDelete(w, r)
			case domain.ActionUndelete.String():
				h.handleUndelete(w, r)
			}
		case common.MIMEApplicationForm:
			switch strings.ToLower(r.FormValue("action")) {
			default:
				h.handleCreate(w, r)
			case domain.ActionDelete.String():
				h.handleDelete(w, r)
			case domain.ActionUndelete.String():
				h.handleUndelete(w, r)
			}
		case common.MIMEMultipartForm:
			h.handleCreate(w, r)
		}
	}
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	req := NewRequestCreate()
	if err := req.bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if r.MultipartForm != nil {
		for k, dst := range map[string]*[]Figure{
			"photo": &req.Properties.Photo,
			"video": &req.Properties.Video,
			"audio": &req.Properties.Audio,
		} {
			file, head, err := r.FormFile(k)
			if err != nil {
				if errors.Is(err, http.ErrMissingFile) {
					continue
				}

				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			location, err := h.media.Upload(r.Context(), domain.File{
				Path:    head.Filename,
				Content: content,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			*dst = append(*dst, Figure{Value: location, Alt: ""})
		}
	}

	in := new(domain.Entry)
	req.populate(in)

	out, err := h.entries.Create(r.Context(), *in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(common.HeaderLocation, out.URL.String())

	links := make([]string, 0)
	for i := range out.Syndications {
		links = append(links, `<`+out.Syndications[i].String()+`>; rel="syndication"`)
	}

	w.Header().Set(common.HeaderLink, strings.Join(links, ", "))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleSource(w http.ResponseWriter, r *http.Request) {
	if r.Method != "" && r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	req := new(RequestSource)
	if err := req.bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	out, err := h.entries.Source(r.Context(), req.URL.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(common.HeaderContentType, common.MIMEApplicationJSONCharsetUTF8)
	if err = json.NewEncoder(w).Encode(NewResponseSource(out, req.Properties...)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	req := new(RequestUpdate)
	if err := req.bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		println(err.Error())

		return
	}

	in, err := h.entries.Source(r.Context(), req.URL.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	req.populate(in)

	out, err := h.entries.Update(r.Context(), req.URL.URL, *in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if out.URL.RequestURI() == req.URL.RequestURI() {
		w.Header().Set(common.HeaderContentType, common.MIMEApplicationJSONCharsetUTF8)
		if err = json.NewEncoder(w).Encode(NewResponseSource(out)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		return
	}

	w.Header().Set(common.HeaderLocation, out.URL.String())
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	req := new(RequestDelete)
	if err := req.bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if _, err := h.entries.Delete(r.Context(), req.URL.URL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(common.HeaderContentType, common.MIMETextPlainCharsetUTF8)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleUndelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	req := new(RequestUndelete)
	if err := req.bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	out, err := h.entries.Undelete(r.Context(), req.URL.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if out.URL.RequestURI() == req.URL.RequestURI() {
		w.Header().Set(common.HeaderContentType, common.MIMEApplicationJSONCharsetUTF8)
		if err = json.NewEncoder(w).Encode(NewResponseSource(out)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		return
	}

	w.Header().Set(common.HeaderLocation, out.URL.String())
	w.WriteHeader(http.StatusCreated)
}

func NewRequestCreate() *RequestCreate {
	return &RequestCreate{
		Type:       make([]string, 0),
		Properties: Properties{},
	}
}

func (r *RequestCreate) bind(req *http.Request) error {
	mediaType, _, err := mime.ParseMediaType(req.Header.Get(common.HeaderContentType))
	if err != nil {
		return fmt.Errorf("cannot understand requested Content-Type: %w", err)
	}

	switch mediaType {
	default:
		return fmt.Errorf("unsupported media type, got '%s', want '%s', '%s' or '%s'", mediaType,
			common.MIMEApplicationJSON, common.MIMEMultipartForm, common.MIMEApplicationForm)
	case common.MIMEApplicationJSON:
		err = json.NewDecoder(req.Body).Decode(r)
	case common.MIMEMultipartForm, common.MIMEApplicationForm:
		in := make(map[string][]string)

		switch mediaType {
		case common.MIMEMultipartForm:
			if err = req.ParseMultipartForm(MaxBodySize); err != nil {
				return fmt.Errorf("cannot parse creation multipart body: %w", err)
			}

			in = req.MultipartForm.Value
		case common.MIMEApplicationForm:
			if err = req.ParseForm(); err != nil {
				return fmt.Errorf("cannot parse creation form body: %w", err)
			}

			in = req.Form
		}

		r.Type = append(r.Type, in["h"]...)

		for k, v := range in {
			switch {
			case strings.HasSuffix(k, "[]"):
				in[strings.TrimSuffix(k, "[]")] = v

				fallthrough
			case strings.HasPrefix(k, "mp-"), k == "h":
				delete(in, k)
			}
		}

		// NOTE(toby3d): hack to encode URL values
		var src []byte
		if src, err = json.Marshal(in); err != nil {
			return fmt.Errorf("cannot marshal creation values for decoding: %w", err)
		}

		err = json.Unmarshal(src, &r.Properties)
	}

	if err != nil {
		return fmt.Errorf("cannot decode creation request body: %w", err)
	}

	for i := range r.Type {
		r.Type[i] = strings.TrimPrefix(r.Type[i], "h-")
	}

	return nil
}

func (r *RequestCreate) populate(dst *domain.Entry) {
	if len(r.Properties.Content) > 0 {
		dst.Content = []byte(r.Properties.Content[0].Value)
	}

	if len(r.Properties.Summary) > 0 {
		dst.Description = r.Properties.Summary[0]
	}

	if len(r.Properties.Published) > 0 {
		dst.PublishedAt = r.Properties.Published[0].Time
	}

	if len(r.Properties.Name) > 0 {
		dst.Title = r.Properties.Name[0]
	}

	if len(r.Properties.Updated) > 0 {
		dst.UpdatedAt = r.Properties.Updated[0].Time
	}

	if len(r.Properties.URL) > 0 {
		dst.URL = r.Properties.URL[0].URL
	}

	dst.Tags = append(dst.Tags, r.Properties.Category...)

	for i := range r.Properties.Photo {
		dst.Photo = append(dst.Photo, r.Properties.Photo[i].Value)
	}

	for i := range r.Properties.Syndication {
		dst.Syndications = append(dst.Syndications, r.Properties.Syndication[i].URL)
	}
}

func (r *RequestSource) bind(req *http.Request) error {
	query := req.URL.Query()
	if r.Q = query.Get("q"); !strings.EqualFold(r.Q, "source") {
		return fmt.Errorf("'q' query MUST be 'source', got '%s'", r.Q)
	}

	var err error
	if r.URL.URL, err = url.Parse(query.Get("url")); err != nil {
		return fmt.Errorf("cannot unmarshal 'url' query: %w", err)
	}

	for k, v := range query {
		if k != "properties" && k != "properties[]" {
			continue
		}

		r.Properties = append(r.Properties, v...)
	}

	return nil
}

func (r *RequestUpdate) bind(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		return fmt.Errorf("cannot decode JSON body: %w", err)
	}

	if !strings.EqualFold(r.Action, "update") {
		return fmt.Errorf("invalid action, got '%s', want '%s'", r.Action, "update")
	}

	return nil
}

func (r RequestUpdate) populate(dst *domain.Entry) {
	if r.Add != nil {
		r.Add.CopyTo(dst)
	}

	if r.Replace != nil {
		if len(r.Replace.Photo) > 0 {
			dst.Photo = make([]*url.URL, 0)
		}

		if len(r.Replace.Category) > 0 {
			dst.Tags = make([]string, 0)
		}

		if len(r.Replace.Syndication) > 0 {
			dst.Syndications = make([]*url.URL, 0)
		}

		r.Replace.CopyTo(dst)
	}

	if r.Delete != nil {
		r.Delete.CopyTo(dst)
	}
}

func (r *RequestDelete) bind(req *http.Request) error {
	mediaType, _, err := mime.ParseMediaType(req.Header.Get(common.HeaderContentType))
	if err != nil {
		return fmt.Errorf("cannot understand requested Content-Type: %w", err)
	}

	switch mediaType {
	default:
		return fmt.Errorf("unsupported media type, got '%s', want '%s' or '%s'", mediaType,
			common.MIMEApplicationJSON, common.MIMEApplicationForm)
	case common.MIMEApplicationJSON:
		err = json.NewDecoder(req.Body).Decode(r)
	case common.MIMEApplicationForm:
		if err = req.ParseForm(); err != nil {
			return fmt.Errorf("cannot decode deletion request: %w", err)
		}

		r.Action = req.PostForm.Get("action")
		if r.URL.URL, err = url.Parse(req.PostForm.Get("url")); err != nil {
			return fmt.Errorf("cannot unmarshal url in deletion request: %w", err)
		}
	}

	if err != nil {
		return fmt.Errorf("cannot parse deletion request: %w", err)
	}

	if !strings.EqualFold(r.Action, "delete") {
		return fmt.Errorf("invalid action, got '%s', want '%s'", r.Action, "delete")
	}

	return nil
}

func (r *RequestUndelete) bind(req *http.Request) error {
	mediaType, _, err := mime.ParseMediaType(req.Header.Get(common.HeaderContentType))
	if err != nil {
		return fmt.Errorf("cannot understand requested Content-Type: %w", err)
	}

	switch mediaType {
	default:
		return fmt.Errorf("unsupported media type, got '%s', want '%s' or '%s'", mediaType,
			common.MIMEApplicationJSON, common.MIMEApplicationForm)
	case common.MIMEApplicationJSON:
		if err = json.NewDecoder(req.Body).Decode(r); err != nil {
			return fmt.Errorf("cannot decode JSON body: %w", err)
		}

		return nil
	case common.MIMEApplicationForm:
		if err = req.ParseForm(); err != nil {
			return fmt.Errorf("cannot parse form body: %w", err)
		}

		r.Action = req.PostFormValue("action")
		if r.URL.URL, err = url.Parse(req.PostFormValue("url")); err != nil {
			return fmt.Errorf("cannot parse url query: %w", err)
		}
	}

	if !strings.EqualFold(r.Action, "undelete") {
		return fmt.Errorf("invalid action, got '%s', want '%s'", r.Action, "undelete")
	}

	return nil
}

func NewResponseSource(src *domain.Entry, properties ...string) *ResponseSource {
	out := &ResponseSource{
		Type: make([]string, 0),
		Properties: Properties{
			Updated:     make([]DateTime, 0),
			Published:   make([]DateTime, 0),
			Photo:       make([]Figure, 0),
			Syndication: make([]URL, 0),
			Content:     make([]Content, 0),
			Category:    make([]string, 0),
			Name:        make([]string, 0),
			Summary:     make([]string, 0),
		},
	}

	if len(properties) == 0 {
		out.Type = append(out.Type, "h-entry")
		properties = []string{
			"updated", "published", "photo", "syndication", "content", "category", "name", "summary",
		}
	}

	if src == nil {
		return out
	}

	for i := range properties {
		switch properties[i] {
		case "updated":
			if src.UpdatedAt.IsZero() {
				continue
			}

			out.Properties.Updated = append(out.Properties.Updated, DateTime{Time: src.UpdatedAt})
		case "published":
			if src.PublishedAt.IsZero() {
				continue
			}

			out.Properties.Published = append(out.Properties.Published, DateTime{
				Time: src.PublishedAt,
			})
		case "photo":
			for j := range src.Photo {
				out.Properties.Photo = append(out.Properties.Photo, Figure{
					Value: src.Photo[j],
				})
			}
		case "syndication":
			for j := range src.Syndications {
				out.Properties.Syndication = append(out.Properties.Syndication, URL{
					URL: src.Syndications[j],
				})
			}
		case "content":
			if len(src.Content) == 0 {
				continue
			}

			out.Properties.Content = append(out.Properties.Content, Content{
				Value: string(src.Content),
			})
		case "category":
			out.Properties.Category = append(out.Properties.Category, src.Tags...)
		case "name":
			out.Properties.Name = append(out.Properties.Name, out.Properties.Name...)
		case "summary":
			if src.Description == "" {
				continue
			}

			out.Properties.Summary = append(out.Properties.Summary, src.Description)
		}
	}

	return out
}

func (p Properties) CopyTo(dst *domain.Entry) {
	if len(p.Updated) > 0 && !p.Updated[0].IsZero() {
		dst.UpdatedAt = p.Updated[0].Time
	}

	if len(p.Published) > 0 && !p.Published[0].IsZero() {
		dst.PublishedAt = p.Published[0].Time
	}

	if len(p.URL) > 0 {
		dst.URL = p.URL[0].URL
	}

	if len(p.Content) > 0 {
		dst.Content = []byte(p.Content[0].Value)
	}

	if len(p.Name) > 0 {
		dst.Title = p.Name[0]
	}

	if len(p.Summary) > 0 {
		dst.Description = p.Summary[0]
	}

	for i := range p.Photo {
		dst.Photo = append(dst.Photo, p.Photo[i].Value)
	}

	dst.Tags = append(dst.Tags, p.Category...)

	for i := range p.Syndication {
		dst.Syndications = append(dst.Syndications, p.Syndication[i].URL)
	}
}

func (f *Figure) UnmarshalJSON(v []byte) error {
	var err error

	buf := new(bufferMedia)

	switch v[0] {
	case '"':
		buf.Value, err = strconv.Unquote(string(v))
	case '{':
		err = json.Unmarshal(v, buf)
	}

	if err != nil {
		return err
	}

	if f.Value, err = url.Parse(buf.Value); err != nil {
		return err
	}

	f.Alt = buf.Alt

	return nil
}

func (f Figure) MarshalJSON() ([]byte, error) {
	if f.Value == nil {
		return []byte(`""`), nil
	}

	return []byte(strconv.Quote(f.Value.String())), nil
}

func (d Delete) CopyTo(dst *domain.Entry) {
	for i := range d.Keys {
		switch d.Keys[i] {
		case "category":
			dst.Tags = make([]string, 0)
		case "content":
			dst.Content = make([]byte, 0)
		case "name":
			dst.Title = ""
		case "photo":
			dst.Photo = make([]*url.URL, 0)
		case "published":
			dst.PublishedAt = time.Time{}
		case "summary":
			dst.Description = ""
		case "updated":
			dst.UpdatedAt = time.Time{}
		case "url":
			dst.URL = new(url.URL)
		case "syndication":
			dst.Syndications = make([]*url.URL, 0)
		}
	}

	// TODO(toby3d): delete property values
}

func (d *Delete) UnmarshalJSON(v []byte) error {
	var err error

	switch v[0] {
	case '[':
		err = json.Unmarshal(v, &d.Keys)
	case '{':
		err = json.Unmarshal(v, &d.Values)
	}

	return err
}

func (u *URL) UnmarshalJSON(v []byte) error {
	raw, err := strconv.Unquote(string(v))
	if err != nil {
		return fmt.Errorf("cannot unqoute URL value: %w", err)
	}

	out, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("cannot parse URL value: %w", err)
	}

	u.URL = out

	return nil
}

func (u URL) MarshalJSON() ([]byte, error) {
	if u.URL == nil {
		return []byte(`""`), nil
	}

	return []byte(strconv.Quote(u.URL.String())), nil
}

func (c Content) String() string {
	if c.HTML == nil {
		return c.Value
	}

	// NOTE(toby3d): trim '<html><head></head><body>' prefix and
	// '</body></html>' suffix
	buf := bytes.NewBuffer(nil)
	if err := html.Render(buf, c.HTML); err != nil {
		return ""
	}

	out := buf.String()

	return out[25 : len(out)-14]
}

func (c *Content) UnmarshalJSON(v []byte) error {
	var err error

	buf := new(bufferHTML)

	switch v[0] {
	case '{':
		err = json.Unmarshal(v, buf)
	case '"':
		c.Value, err = strconv.Unquote(string(v))
	}

	if err != nil {
		return err
	}

	if buf.HTML == "" {
		return nil
	}

	if c.HTML, err = html.Parse(strings.NewReader(buf.HTML)); err != nil {
		return err
	}

	return nil
}

func (c Content) MarshalJSON() ([]byte, error) {
	if c.HTML == nil {
		return []byte(strconv.Quote(c.Value)), nil
	}

	buf := bytes.NewBuffer(nil)
	if err := html.Render(buf, c.HTML); err != nil {
		return nil, err
	}

	out := buf.String()

	// NOTE(toby3d): trim '<html><head></head><body>' prefix and
	// '</body></html>' suffix
	return []byte(`{"html":"` + out[25:len(out)-14] + `"}`), nil
}

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	v, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	var out time.Time

	for _, format := range []string{
		time.RFC3339,
		// NOTE(toby3d): fallback for datetime-local input HTML node
		// format
		"2006-01-02T15:04",
	} {
		if out, err = time.Parse(format, v); err != nil {
			continue
		}

		dt.Time = out

		return nil
	}

	return err
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return nil, nil
	}

	return []byte(strconv.Quote(dt.Format(time.RFC3339))), nil
}

func (a *Action) UnmarshalJSON(b []byte) error {
	v, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	out, err := domain.ParseAction(strings.TrimSpace(strings.ToLower(v)))
	if err != nil {
		return err
	}

	a.Action = out

	return nil
}

func (a Action) MarshalJSON() ([]byte, error) {
	if a.Action == domain.ActionUnd {
		return []byte(`""`), nil
	}

	return []byte(strconv.Quote(a.Action.String())), nil
}
