package simple_router

import (
	"github.com/go-chi/chi/v5"
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
}

type Config struct {
	Port string `json:"port"`
	Name string `json:"name"`
}
