package modules

import (
	"context"
	"errors"
	"fmt"
	"github.com/advancedlogic/goms/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	middleware     []func(ctx *gin.Context)
	websiteFolder  map[string]string
	cert           string
	key            string
	logger         *logrus.Logger
	server         *http.Server
}

// Rest builder
type RestBuilder struct {
	*models.Environment
	*Rest
	models.Exception
}

func NewRestBuilder(environment *models.Environment) *RestBuilder {
	rb := &RestBuilder{
		Environment: environment,
		Rest: &Rest{
			port:           environment.GetIntOrDefault("transport.port", 8080),
			readTimeout:    environment.GetDurationOrDefault("transport.readTimeout", 10) * time.Second,
			writeTimeout:   environment.GetDurationOrDefault("transport.writeTimeout", 10) * time.Second,
			getHandlers:    make(map[string][]gin.HandlerFunc),
			postHandlers:   make(map[string][]gin.HandlerFunc),
			putHandlers:    make(map[string][]gin.HandlerFunc),
			deleteHandlers: make(map[string][]gin.HandlerFunc),
			middleware:     make([]func(ctx *gin.Context), 0),
			websiteFolder:  make(map[string]string),
			logger:         environment.Logger,
		},
	}

	return rb.
		WithPort(environment.GetIntOrDefault("service.port", 8080)).
		WithReadTimeout(environment.GetDurationOrDefault("service.timeout", 10*time.Second)).
		WithWriteTimeout(environment.GetDurationOrDefault("service.timeout", 10*time.Second)).
		WithStaticFilesFolder("/static", environment.GetStringOrDefault("service.public", ""))

}

func (rb *RestBuilder) WithPort(port int) *RestBuilder {
	if port == 0 {
		rb.Catch("port must be greater than zero")
	}
	rb.port = port
	return rb
}

func (rb *RestBuilder) WithReadTimeout(duration time.Duration) *RestBuilder {
	if duration == 0 {
		rb.Catch("read timeout cannot be zero")
	}
	rb.readTimeout = duration
	return rb
}

func (rb *RestBuilder) WithWriteTimeout(duration time.Duration) *RestBuilder {
	if duration == 0 {
		rb.Catch("read timeout cannot be zero")
	}
	rb.writeTimeout = duration
	return rb
}

func (rb *RestBuilder) WithGetHandler(path string, handler func(ctx *gin.Context)) *RestBuilder {
	if path == "" {
		rb.Catch("path cannot be empty")
	}
	if _, exists := rb.getHandlers[path]; !exists {
		rb.getHandlers[path] = make([]gin.HandlerFunc, 0)
	}
	rb.getHandlers[path] = append(rb.getHandlers[path], gin.HandlerFunc(handler))
	return rb
}

func (rb *RestBuilder) WithPostHandler(path string, handler func(ctx *gin.Context)) *RestBuilder {
	if path == "" {
		rb.Catch("path cannot be empty")
	}
	if _, exists := rb.postHandlers[path]; !exists {
		rb.postHandlers[path] = make([]gin.HandlerFunc, 0)
	}
	rb.postHandlers[path] = append(rb.postHandlers[path], gin.HandlerFunc(handler))
	return rb
}

func (rb *RestBuilder) WithPutHandler(path string, handler func(ctx *gin.Context)) *RestBuilder {
	if path == "" {
		rb.Catch("path cannot be empty")
	}
	if _, exists := rb.putHandlers[path]; !exists {
		rb.putHandlers[path] = make([]gin.HandlerFunc, 0)
	}
	rb.putHandlers[path] = append(rb.putHandlers[path], gin.HandlerFunc(handler))
	return rb
}

func (rb *RestBuilder) WithDeleteHandler(path string, handler func(ctx *gin.Context)) *RestBuilder {
	if path == "" {
		rb.Catch("path cannot be empty")
	}
	if _, exists := rb.deleteHandlers[path]; !exists {
		rb.deleteHandlers[path] = make([]gin.HandlerFunc, 0)
	}
	rb.deleteHandlers[path] = append(rb.deleteHandlers[path], gin.HandlerFunc(handler))
	return rb
}

func (rb *RestBuilder) WithMiddleware(middleware func(ctx *gin.Context)) *RestBuilder {
	rb.middleware = append(rb.middleware, middleware)
	return rb
}

func (rb *RestBuilder) WithStaticFilesFolder(uri, folder string) *RestBuilder {
	if uri != "" && folder != "" {
		rb.websiteFolder[uri] = folder
	}
	return rb
}

func (rb *RestBuilder) WithTLS(pem, key string) *RestBuilder {
	if pem == "" && key == "" {
		rb.Catch("certificate or key cannot be empty")
	}
	rb.cert = pem
	rb.key = key
	return rb
}

func (rb *RestBuilder) Build() (*Rest, error) {
	ginMode := rb.GetStringOrDefault("service.mode", "release")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := rb.CheckErrors(rb.Errors()); err != nil {
		return nil, err
	}
	if rb.Rest != nil {
		return rb.Rest, nil
	}
	return nil, errors.New("")
}

func (r *Rest) GetHandler(endpoint string, handler interface{}) {
	if _, exists := r.getHandlers[endpoint]; !exists {
		r.getHandlers[endpoint] = make([]gin.HandlerFunc, 0)
	}
	r.getHandlers[endpoint] = append(r.getHandlers[endpoint], gin.HandlerFunc(handler.(func(*gin.Context))))
}

func (r *Rest) PostHandler(endpoint string, handler interface{}) {
	if _, exists := r.postHandlers[endpoint]; !exists {
		r.postHandlers[endpoint] = make([]gin.HandlerFunc, 0)
	}
	r.postHandlers[endpoint] = append(r.postHandlers[endpoint], gin.HandlerFunc(handler.(func(ctx *gin.Context))))
}

func (r *Rest) PutHandler(endpoint string, handler interface{}) {
	if _, exists := r.putHandlers[endpoint]; !exists {
		r.putHandlers[endpoint] = make([]gin.HandlerFunc, 0)
	}
	r.putHandlers[endpoint] = append(r.putHandlers[endpoint], gin.HandlerFunc(handler.(func(ctx *gin.Context))))
}

func (r *Rest) DeleteHandler(endpoint string, handler interface{}) {
	if _, exists := r.deleteHandlers[endpoint]; !exists {
		r.deleteHandlers[endpoint] = make([]gin.HandlerFunc, 0)
	}
	r.deleteHandlers[endpoint] = append(r.deleteHandlers[endpoint], gin.HandlerFunc(handler.(func(ctx *gin.Context))))
}

func (r *Rest) Middleware(middleware interface{}) {
	if r.middleware == nil {
		r.middleware = make([]func(*gin.Context), 0)
	}
	r.middleware = append(r.middleware, middleware.(func(ctx *gin.Context)))
}

func (r *Rest) StaticFilesFolder(uri, folder string) {
	r.websiteFolder[uri] = folder
}

func (r *Rest) Run() error {
	router := gin.New()
	router.Use(ginlogrus.Logger(r.logger), gin.Recovery())
	router.GET("/healthcheck", func(c *gin.Context) {
		c.String(200, "product service is good")
	})

	for _, middleware := range r.middleware {
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
