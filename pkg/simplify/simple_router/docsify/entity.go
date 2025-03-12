package docsify

import "net/http"

type Service interface {
	Apply() http.HandlerFunc
}
