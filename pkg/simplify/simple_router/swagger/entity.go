package swagger

import "net/http"

type Service interface {
	Apply() http.HandlerFunc
}
