package http

import (
	"bytes"
	"encoding/json"
	"path"

	"github.com/fasthttp/router"
	"github.com/lestrrat-go/jwx/jwa"
	http "github.com/valyala/fasthttp"
	"golang.org/x/xerrors"

	"source.toby3d.me/toby3d/middleware"
	"source.toby3d.me/website/micropub/internal/common"
	"source.toby3d.me/website/micropub/internal/domain"
	"source.toby3d.me/website/micropub/internal/media"
)

// RequestHandler represents a handler with business logic for HTTP requests.
type RequestHandler struct {
	useCase media.UseCase
}

// New creates a new HTTP delivery handler.
func New(useCase media.UseCase) *RequestHandler {
	return &RequestHandler{
		useCase: useCase,
	}
}

// Register register media endpoints for router.
func (h *RequestHandler) Register(r *router.Router) {
	chain := middleware.Chain{
		middleware.JWTWithConfig(middleware.JWTConfig{
			AuthScheme:    "Bearer",
			ContextKey:    "token",
			SigningKey:    []byte("hackme"),                // TODO(toby3d): replace setting from config
			SigningMethod: jwa.SignatureAlgorithm("HS256"), // TODO(toby3d): replace setting from config
			TokenLookup:   middleware.SourceHeader + ":" + http.HeaderAuthorization,
		}),
	}
	// TODO(toby3d): The Media Endpoint MUST accept the same access tokens
	// that the Micropub endpoint accepts.
	r.POST("/media", chain.RequestHandler(h.Update))
	r.GET("/media/{fileName:*}", chain.RequestHandler(h.Read))
}

func (h *RequestHandler) Update(ctx *http.RequestCtx) {
	ctx.SetContentType(common.MIMEApplicationJSON)
	encoder := json.NewEncoder(ctx)

	ff, err := ctx.FormFile("file")
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		encoder.Encode(&domain.Error{
			Code:        "invalid_request",
			Description: err.Error(),
			Frame:       xerrors.Caller(1),
		})

		return
	}

	f, err := ff.Open()
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		encoder.Encode(&domain.Error{
			Code:        "invalid_request",
			Description: err.Error(),
			Frame:       xerrors.Caller(1),
		})

		return
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err = buf.ReadFrom(f); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		encoder.Encode(&domain.Error{
			Code:        "invalid_request",
			Description: err.Error(),
			Frame:       xerrors.Caller(1),
		})

		return
	}

	fileName, err := h.useCase.Upload(ctx, ff.Filename, buf.Bytes())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		encoder.Encode(&domain.Error{
			Code:        "invalid_request",
			Description: err.Error(),
			Frame:       xerrors.Caller(1),
		})

		return
	}

	ctx.SetStatusCode(http.StatusCreated)
	ctx.Response.Header.Set(http.HeaderLocation, path.Join("/", "media", fileName))
	encoder.Encode(struct{}{})
}

func (h *RequestHandler) Read(ctx *http.RequestCtx) {
	ctx.SetContentType(common.MIMEApplicationJSON)
	encoder := json.NewEncoder(ctx)

	fileName, ok := ctx.UserValue("fileName").(string)
	if !ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		encoder.Encode(&domain.Error{
			Code:        "invalid_request",
			Description: "media file name is not provided",
			Frame:       xerrors.Caller(1),
		})

		return
	}

	contents, err := h.useCase.Download(ctx, fileName)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		encoder.Encode(&domain.Error{
			Code:        "invalid_request",
			Description: err.Error(),
			Frame:       xerrors.Caller(1),
		})

		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType(common.MIMEApplicationOctetStream)
	ctx.SetBody(contents)
}
