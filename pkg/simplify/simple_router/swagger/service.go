package swagger

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
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
