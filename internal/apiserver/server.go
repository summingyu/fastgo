package apiserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/onexstack/fastgo/internal/apiserver/biz"
	"github.com/onexstack/fastgo/internal/apiserver/handler"
	"github.com/onexstack/fastgo/internal/apiserver/pkg/validation"
	"github.com/onexstack/fastgo/internal/apiserver/store"
	"github.com/onexstack/fastgo/internal/pkg/core"
	"github.com/onexstack/fastgo/internal/pkg/errorsx"
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

	// 初始化数据库连接
	db, err := cfg.MySQLOptions.NewDB()
	if err != nil {
		return nil, err
	}
	store := store.NewStore(db)
	cfg.InstallRESTAPI(engine, store)

	httpsrv := &http.Server{Addr: cfg.Addr, Handler: engine}
	return &Server{cfg: cfg, srv: httpsrv}, nil
}

func (cfg *Config) InstallRESTAPI(engine *gin.Engine, store store.IStore) {
	// 注册404 Handler
	engine.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errorsx.ErrNotFound.WithMessage("Page not found"), nil)
	})
	// 注册/healthz handler.
	engine.GET("/healthz", func(c *gin.Context) {
		core.WriteResponse(c, map[string]string{"status": "ok"}, nil)
	})

	handler := handler.NewHandler(biz.NewBiz(store), validation.NewValidator(store))
	authMiddlewares := []gin.HandlerFunc{}

	v1 := engine.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", handler.CreateUser)
			userv1.PUT(":userID", handler.UpdateUser)
			userv1.DELETE(":userID", handler.DeleteUser)
			userv1.GET(":userID", handler.GetUser)
			userv1.GET("", handler.ListUser)
		}

		postv1 := v1.Group("/posts", authMiddlewares...)
		{
			postv1.POST("", handler.CreatePost)
			postv1.PUT(":postID", handler.UpdatePost)
			postv1.DELETE(":postID", handler.DeletePost)
			postv1.GET(":postID", handler.GetPost)
			postv1.GET("", handler.ListPost)
		}
	}
}

func (s *Server) Run() error {
	slog.Info("Start to listening to incoming requests on http address", "addr", s.cfg.Addr)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置一个10秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
	// 创建一个10秒的超时上下文，然后调用Shutdown方法来优雅地关闭服务器。
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 优雅关闭服务器
	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("Insecure Server forced to shutdown", "err", err)
		return err
	}
	slog.Info("Server exited normally.")
	return nil
}
