package main

import (
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env/v9"

	"source.toby3d.me/toby3d/pub/internal/domain"
	mediahttpdelivery "source.toby3d.me/toby3d/pub/internal/media/delivery/http"
	mediamemoryrepo "source.toby3d.me/toby3d/pub/internal/media/repository/memory"
	mediaucase "source.toby3d.me/toby3d/pub/internal/media/usecase"
	"source.toby3d.me/toby3d/pub/internal/urlutil"
)

var (
	config = new(domain.Config)
	logger = log.New(os.Stdout, "Micropub	", log.LstdFlags)
)

func init() {
	if err := env.ParseWithOptions(config, env.Options{
		Prefix: "MICROPUB_",
	}); err != nil {
		logger.Fatal("cannot parse environment variables into config:", err)
	}
}

func main() {
	mediaRepo := mediamemoryrepo.NewMemoryMediaRepository()
	mediaUseCase := mediaucase.NewMediaUseCase(mediaRepo)
	mediaHandler := mediahttpdelivery.NewHandler(mediaUseCase, *config)

	server := http.Server{
		Addr: config.HTTP.Bind,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			head, _ := urlutil.ShiftPath(r.RequestURI)

			switch head {
			default:
				http.NotFound(w, r)
			case "media":
				mediaHandler.ServeHTTP(w, r)
			}
		}),
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("cannot listen and serve:", err)
	}
}
