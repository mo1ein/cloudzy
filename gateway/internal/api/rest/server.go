package rest

import (
	"context"
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Server struct {
	engine *gin.Engine
}

func New() *Server {
	//if appEnv == config.ProductionEnv {
	//	gin.SetMode(gin.ReleaseMode)
	//}

	r := gin.New()

	//err := r.SetTrustedProxies(trustedProxies)
	//if err != nil {
	//	log.Fatal("trusted proxies has error", log.J{
	//		"error": err.Error(),
	//	})
	//}

	r.RedirectTrailingSlash = false

	r.Use(
		requestid.New(
			requestid.WithGenerator(func() string {
				return ""
			}),
			requestid.WithCustomHeaderStrKey("x-request-id"),
		),
	)

	return &Server{
		engine: r,
	}
}

const headerTimeout = 10 * time.Second

func (s *Server) Serve(ctx context.Context, address string) error {
	srv := &http.Server{
		Addr:              address,
		Handler:           s.engine,
		ReadHeaderTimeout: headerTimeout,
	}

	// todo: do this with log.info
	log.Println(fmt.Sprintf("rest server starting at: %s", address))
	srvError := make(chan error)
	go func() {
		srvError <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		// todo: do this with log.info
		log.Println("rest server is shutting down")
		return srv.Shutdown(ctx)
	case err := <-srvError:
		return err
	}
}
