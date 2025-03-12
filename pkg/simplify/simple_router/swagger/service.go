package swagger

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type service struct {
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

func (s *service) Apply() http.HandlerFunc {
	return httpSwagger.WrapHandler.ServeHTTP
}
