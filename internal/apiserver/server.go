package apiserver

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	mw "github.com/onexstack/fastgo/internal/pkg/middleware"
	genericoptions "github.com/onexstack/fastgo/pkg/options"
)

type Config struct {
	MySQLOptions *genericoptions.MySQLOptions
	Addr         string
}

type Server struct {
	cfg *Config
	srv *http.Server
}

func (cfg *Config) NewServer() (*Server, error) {
	engine := gin.New()
	// gin.Recovery() 中间件，用来捕获任何 panic，并恢复
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.RequestID()}
	engine.Use(mws...)
	// 注册404 Handler
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PageNotFound", "message": "Page not found."})
	})
	// 注册/healthz handler.
	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	httpsrv := &http.Server{Addr: cfg.Addr, Handler: engine}
	return &Server{cfg: cfg, srv: httpsrv}, nil
}

func (s *Server) Run() error {
	slog.Info("Read MySQL host from config", "mysql.addr", s.cfg.MySQLOptions.Addr)
	slog.Info("Start to listening to incoming requests on http address", "addr", s.cfg.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
