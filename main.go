//go:generate go install github.com/valyala/quicktemplate/qtc@master
//go:generate qtc -dir=web/template
//go:generate go install golang.org/x/text/cmd/gotext@master
//go:generate gotext -srclang=en update -lang=en,ru -out=locales_gen.go
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/caarlos0/env/v9"

	"source.toby3d.me/toby3d/pub/internal/domain"
	mediahttpdelivery "source.toby3d.me/toby3d/pub/internal/media/delivery/http"
	mediamemoryrepo "source.toby3d.me/toby3d/pub/internal/media/repository/memory"
	mediaucase "source.toby3d.me/toby3d/pub/internal/media/usecase"
	"source.toby3d.me/toby3d/pub/internal/urlutil"
)

var (
	config = new(domain.Config)
	logger = log.New(os.Stdout, "Micropub\t", log.LstdFlags)
)

var cpuProfilePath, memProfilePath string

func init() {
	flag.StringVar(&cpuProfilePath, "cpuprofile", "", "set path to saving CPU memory profile")
	flag.StringVar(&memProfilePath, "memprofile", "", "set path to saving pprof memory profile")
	flag.Parse()

	if err := env.ParseWithOptions(config, env.Options{
		Prefix: "MICROPUB_",
	}); err != nil {
		logger.Fatal("cannot parse environment variables into config:", err)
	}
}

func main() {
	ctx := context.Background()
	mediaRepo := mediamemoryrepo.NewMemoryMediaRepository()
	mediaUseCase := mediaucase.NewMediaUseCase(mediaRepo)
	mediaHandler := mediahttpdelivery.NewHandler(mediaUseCase, *config)

	server := http.Server{
		ErrorLog: logger,
		Addr:     config.HTTP.Bind,
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

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if cpuProfilePath != "" {
		cpuProfile, err := os.Create(cpuProfilePath)
		if err != nil {
			logger.Fatalln("could not create CPU profile:", err)
		}
		defer cpuProfile.Close()

		if err = pprof.StartCPUProfile(cpuProfile); err != nil {
			logger.Fatalln("could not start CPU profile:", err)
		}
		defer pprof.StopCPUProfile()
	}

	go func() {
		logger.Printf("started at %s, available at %s", config.HTTP.Bind, config.HTTP.BaseURL())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalln("cannot listen and serve:", err)
		}
	}()

	<-done

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalln("failed shutdown of server:", err)
	}

	if memProfilePath == "" {
		return
	}

	memProfile, err := os.Create(memProfilePath)
	if err != nil {
		logger.Fatalln("could not create memory profile:", err)
	}
	defer memProfile.Close()

	runtime.GC() // NOTE(toby3d): get up-to-date statistics

	if err = pprof.WriteHeapProfile(memProfile); err != nil {
		logger.Fatalln("could not write memory profile:", err)
	}
}
