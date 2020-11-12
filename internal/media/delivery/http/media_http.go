package http

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"

	"source.toby3d.me/toby3d/pub/internal/common"
	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/media"
)

type (
	Handler struct {
		media  media.UseCase
		config domain.Config
	}

	Error struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description,omitempty"`
	}
)

func NewHandler(media media.UseCase, config domain.Config) *Handler {
	return &Handler{
		media:  media,
		config: config,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, "method MUST be "+http.MethodPost, http.StatusBadRequest)

		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get(common.HeaderContentType))
	if err != nil {
		WriteError(w, "Content-Type header MUST be "+common.MIMEMultipartForm, http.StatusBadRequest)

		return
	}

	if mediaType != common.MIMEMultipartForm {
		WriteError(w, "Content-Type header MUST be "+common.MIMEMultipartForm, http.StatusBadRequest)

		return
	}

	file, head, err := r.FormFile("file")
	if err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)

		return
	}
	defer file.Close()

	in := &domain.File{
		Path:    head.Filename,
		Content: make([]byte, 0),
	}

	if in.Content, err = io.ReadAll(file); err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)

		return
	}

	out, err := h.media.Upload(r.Context(), *in)
	if err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set(common.HeaderLocation, h.config.HTTP.BaseURL().JoinPath(out.Path).String())
	w.WriteHeader(http.StatusCreated)
}

func WriteError(w http.ResponseWriter, description string, status int) {
	out := &Error{ErrorDescription: description}

	switch status {
	case http.StatusBadRequest:
		out.Error = "invalid_request"
	case http.StatusForbidden: // TODO(toby3d): insufficient_scope
		out.Error = "forbidden"
	case http.StatusUnauthorized:
		out.Error = "unauthorized"
	}

	_ = json.NewEncoder(w).Encode(out)

	w.Header().Set(common.HeaderContentType, common.MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(status)
}
