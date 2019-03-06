package modules

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/toorop/gin-logrus"

	"net/http"
	"time"
)

// Rest Server
type Rest struct {
	// Port bound to server
	port           int
	readTimeout    time.Duration
	writeTimeout   time.Duration
	getHandlers    map[string][]gin.HandlerFunc
	postHandlers   map[string][]gin.HandlerFunc
	putHandlers    map[string][]gin.HandlerFunc
	deleteHandlers map[string][]gin.HandlerFunc
	middlewares    []func(ctx *gin.Context)
	websiteFolder  map[string]string
	cert           string
	key            string
	server         *http.Server
}

// Rest builder
type RestBuilder struct {
	*Environment
	*Rest
}

func NewRestBuilder(environment *Environment) *RestBuilder {
	return &RestBuilder{
		Environment: environment,
		Rest: &Rest{
			port:           8080,
			readTimeout:    10 * time.Second,
			writeTimeout:   10 * time.Second,
			getHandlers:    make(map[string][]gin.HandlerFunc),
			postHandlers:   make(map[string][]gin.HandlerFunc),
			putHandlers:    make(map[string][]gin.HandlerFunc),
			deleteHandlers: make(map[string][]gin.HandlerFunc),
			middlewares:    make([]func(ctx *gin.Context), 0),
			websiteFolder:  make(map[string]string),
		},
	}
}

func (rb *RestBuilder) WithPort(port int) *RestBuilder {
	rb.port = port
	return rb
}

func (rb *RestBuilder) WithReadTimeout(duration time.Duration) *RestBuilder {
	rb.readTimeout = duration
	return rb
}

func (rb *RestBuilder) WithWriteTimeout(duration time.Duration) *RestBuilder {
	rb.writeTimeout = duration
	return rb
}

func (rb *RestBuilder) WithGetHandler(path string, handler func(ctx *gin.Context)) *RestBuilder {
	if _, exists := rb.getHandlers[path]; !exists {
		rb.getHandlers[path] = make([]gin.HandlerFunc, 0)
	}
	rb.getHandlers[path] = append(rb.getHandlers[path], gin.HandlerFunc(handler))
	return rb
}

func (rb *RestBuilder) WithPostHandler(path string, handler func(ctx *gin.Context)) *RestBuilder {
	if _, exists := rb.postHandlers[path]; !exists {
		rb.postHandlers[path] = make([]gin.HandlerFunc, 0)
	}
	rb.postHandlers[path] = append(rb.postHandlers[path], gin.HandlerFunc(handler))
	return rb
}

func (rb *RestBuilder) WithPutHandler(path string, handler func(ctx *gin.Context)) *RestBuilder {
	if _, exists := rb.putHandlers[path]; !exists {
		rb.putHandlers[path] = make([]gin.HandlerFunc, 0)
	}
	rb.putHandlers[path] = append(rb.putHandlers[path], gin.HandlerFunc(handler))
	return rb
}

func (rb *RestBuilder) WithDeleteHandler(path string, handler func(ctx *gin.Context)) *RestBuilder {
	if _, exists := rb.deleteHandlers[path]; !exists {
		rb.deleteHandlers[path] = make([]gin.HandlerFunc, 0)
	}
	rb.deleteHandlers[path] = append(rb.deleteHandlers[path], gin.HandlerFunc(handler))
	return rb
}

func (rb *RestBuilder) WithMiddleware(middleware func(ctx *gin.Context)) *RestBuilder {
	rb.middlewares = append(rb.middlewares, middleware)
	return rb
}

func (rb *RestBuilder) WithStaticFilesFolder(uri, folder string) *RestBuilder {
	rb.websiteFolder[uri] = folder
	return rb
}

func (rb *RestBuilder) WithTLS(pem, key string) *RestBuilder {
	rb.cert = pem
	rb.key = key
	return rb
}

func (rb *RestBuilder) Build() (*Rest, error) {
	if rb.Rest != nil {
		return rb.Rest, nil
	}
	return nil, errors.New("")
}

func (r *Rest) Run(opts ...string) error {
	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	for _, middleware := range r.middlewares {
		router.Use(gin.HandlerFunc(middleware))
	}

	for path, handlers := range r.getHandlers {
		router.GET(path, handlers...)
	}

	for path, handler := range r.postHandlers {
		router.POST(path, handler...)
	}
	for path, handler := range r.putHandlers {
		router.PUT(path, handler...)
	}

	for path, handler := range r.deleteHandlers {
		router.DELETE(path, handler...)
	}

	if len(r.websiteFolder) > 0 {
		for uri, folder := range r.websiteFolder {
			router.Static(uri, folder)
		}
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", r.port),
		Handler:        router,
		ReadTimeout:    r.readTimeout,
		WriteTimeout:   r.writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	if r.cert != "" && r.key != "" {
		if err := s.ListenAndServeTLS(r.cert, r.key); err != nil {
			return err
		}
	} else if err := s.ListenAndServe(); err != nil {
		return err
	}

	r.server = s
	return nil
}

func (r *Rest) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.server.Shutdown(ctx)
}
