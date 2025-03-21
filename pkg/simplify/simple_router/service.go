package simple_router

import (
	"net/http"
	"net/http/pprof"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/simple_router/docsify"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/simple_router/ping"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/app_profile"
)

var _ Service = (*App)(nil)

func NewService(c Config) *App {
	routes := initRoutes()
	return &App{
		Router: routes,
		Port:   setPort(c.Port),
	}
}

func (a *App) Run() error {
	return http.ListenAndServe(a.Port, a.Router)
}

func registerPprofRoutes(router chi.Router) {
	pprofMux := http.NewServeMux()
	pprofMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	pprofMux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	pprofMux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	pprofMux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	pprofMux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	router.Mount("/debug", pprofMux)
}

func initRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", ping.NewService().Apply())

	if !app_profile.IsProdProfile() {
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		r.Handle("/documentation-tech/*", docsify.NewService().Apply())
		registerPprofRoutes(r)
	}
	return r
}

func setPort(p string) string {
	if p != "" {
		return ":" + p
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	return ":" + appDefaultPort
}
