package simple_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

const (
	appDefaultPort = "8080"
)

type Service interface {
	Run() error
	RegisterRoute(pattern string, handler http.HandlerFunc)
}

type App struct {
	Router *chi.Mux
	Port   string
	log    log.Service
}

type Config struct {
	Port string
	Name string
}
