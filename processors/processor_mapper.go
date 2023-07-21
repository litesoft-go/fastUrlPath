package processors

import (
	"net/http"
	"net/url"
)

type ProcessorMapper interface {
	FromRequest(req *http.Request) (Processor, error)
	FromUrl(url *url.URL) (Processor, error)
	FromPath(path string) (Processor, error)
}
