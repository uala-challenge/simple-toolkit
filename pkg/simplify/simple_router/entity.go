package simple_router

import (
	"github.com/go-chi/chi/v5"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

const (
	appDefaultPort = "8080"
)

type Service interface {
	Run() error
}

type App struct {
	Router *chi.Mux
	Port   string
	log    log.Service
}

type Config struct {
	Port string `json:"port"`
	Name string `json:"name"`
}
