package router

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Endpoint struct {
	mux        *http.ServeMux
	middleware []Adapter
	paths      []string
}

func initApi(mux *http.ServeMux, path string, middleware ...Adapter) *Endpoint {
	return &Endpoint{
		mux:        mux,
		middleware: middleware,
		paths:      []string{path},
	}
}
func (e *Endpoint) group(path string, middleware ...Adapter) *Endpoint {
	ep := new(Endpoint)
	ep.mux = e.mux
	ep.middleware = append(e.middleware, middleware...)
	ep.paths = append(e.paths, path)

	return ep
}
func (e *Endpoint) endpoint(path ...string) string {
	joinPath, err := url.JoinPath("/", append(e.paths, path...)...)

	if err != nil {
		log.Fatal().Err(err).Msg("endpoint error")
	}

	return joinPath
}

func (ph PathHandler) log(pattern string) {
	v := reflect.ValueOf(ph)
	methods := make([]string, 0, 5)

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		if f.Kind() == reflect.Func && !f.IsNil() {
			methods = append(methods, v.Type().Field(i).Name)
		}
	}

	log.Info().Msgf("endpoint: %-40s methods: %s", pattern, strings.Join(methods, " "))
}

func (e *Endpoint) handle(h ...PathHandler) {
	for _, handler := range h {
		pattern := e.endpoint(handler.Path)
		handler.log(pattern)

		if e.middleware != nil {
			e.mux.Handle(pattern, Adapt(handler, e.middleware...))
		} else {
			e.mux.Handle(pattern, handler)
		}
	}
}

type PathHandler struct {
	Path                          string
	GET, POST, PUT, PATCH, DELETE http.HandlerFunc
}

func (ph PathHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		if ph.GET != nil {
			ph.GET(w, r)
			return
		}

		fallthrough

	case "POST":
		if ph.POST != nil {
			ph.POST(w, r)
			return
		}

		fallthrough

	case "PUT":
		if ph.PUT != nil {
			ph.PUT(w, r)
			return
		}

		fallthrough

	case "PATCH":
		if ph.PATCH != nil {
			ph.PATCH(w, r)
			return
		}

		fallthrough

	case "DELETE":
		if ph.DELETE != nil {
			ph.DELETE(w, r)
			return
		}

		fallthrough

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
