package docsify

import (
	"net/http"
)

type service struct {
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

// Apply godoc
// @Summary      Serve static documentation files
// @Description  Serves Docsify documentation files from the /docs directory
// @Tags         documentation
// @Accept       json
// @Produce      json
// @Success      200  {string}  string "OK"
// @Failure      404  {string}  string "Not Found"
// @Router       /docs/* [get]
func (s *service) Apply() http.HandlerFunc {
	return http.StripPrefix("/docs", http.FileServer(http.Dir("./docs"))).ServeHTTP
}
