package main

import (
	"context"
	"github.com/ValeryVerkhoturov/chat/auth"
	"github.com/ValeryVerkhoturov/chat/config"
	"github.com/ValeryVerkhoturov/chat/db"
	"github.com/ValeryVerkhoturov/chat/handlers"
	v1Handlers "github.com/ValeryVerkhoturov/chat/handlers/v1"
	v1Socket "github.com/ValeryVerkhoturov/chat/handlers/v1/socket"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func port() string {
	port := config.Port
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

func createCorsMiddleware() gin.HandlerFunc {
	var corsConfig = cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Origins
	return cors.New(corsConfig)
}

func createEngine() *graceful.Graceful {
	router, err := graceful.Default()
	if err != nil {
		panic(err)
	}

	router.Use(createCorsMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(sessions.Sessions("chat-session", cookie.NewStore([]byte(config.SessionSecret))))

	router.Static("/images/", "./public/images")
	router.StaticFile("/css/output.css", "./public/css/output.css")
	router.LoadHTMLGlob("templates/templates/*")

	router.GET("/", handlers.Index)
	router.GET("/index.html", handlers.Index)

	v1Router := router.Group("/v1")
	v1Router.Use(auth.CreateSessionIfNotExists)
	{
		v1Router.GET("/chat-widget", v1Handlers.ChatWidget)
	}
	v1Router.Use(auth.SessionRequired)
	{
		v1Router.GET("/chat-container", v1Handlers.ChatContainer)
	}
	var hub = v1Socket.NewHub()
	go hub.Run()
	v1Router.GET("/ws", func(c *gin.Context) {
		v1Socket.WS(c, hub)
	})

	return router
}

func main() {
	var err error

	// Log init
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// MongoDB connect
	mongoCtx, cancel := db.InitDB()
	defer cancel()
	defer func() {
		if err = db.MongoClient.Disconnect(mongoCtx); err != nil {
			panic(err)
		}
	}()

	// Graceful termination when shutting down a process init
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Gin init
	log.Info("Starting Server on http://localhost" + port())
	router := createEngine()
	defer router.Close()

	if err = router.RunWithContext(ctx); err != nil && err != context.Canceled {
		log.Fatal("Unable to start:", err)
	}
}
