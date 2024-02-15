package main

import (
	"context"
	"github.com/ValeryVerkhoturov/chat/auth"
	"github.com/ValeryVerkhoturov/chat/config"
	"github.com/ValeryVerkhoturov/chat/db"
	"github.com/ValeryVerkhoturov/chat/handlers"
	v1Handlers "github.com/ValeryVerkhoturov/chat/handlers/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

func engine() *graceful.Graceful {
	router, err := graceful.Default()
	if err != nil {
		panic(err)
	}

	router.Use(cors.Default()) // AllowAllOrigins true
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(sessions.Sessions("chat-session", cookie.NewStore([]byte(config.SessionSecret))))

	router.Static("/images/", "./public/images")
	router.StaticFile("/css/output-css", "./public/css/output.css")
	router.LoadHTMLGlob("templates/templates/*")

	router.GET("/", handlers.Index)
	router.GET("/index.html", handlers.Index)

	v1 := router.Group("/v1")
	v1.Use(auth.CreateSessionIfNotExists)
	{
		v1.GET("/chat-widget", v1Handlers.ChatWidget)
	}
	v1.Use(auth.SessionRequired)
	{
		v1.GET("/chat-container", v1Handlers.ChatContainer)
	}
	return router
}

func main() {
	var err error

	// Log init
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Graceful termination when shutting down a process init
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// MongoDB connect
	mongoCtx, cancel := db.InitDB()
	defer cancel()
	defer func() {
		if err = db.MongoClient.Disconnect(mongoCtx); err != nil {
			panic(err)
		}
	}()

	// Gin init
	log.Info("Starting Server on http://localhost" + port())
	router := engine()
	defer router.Close()
	if err = router.RunWithContext(ctx); err != nil && err != context.Canceled {
		log.Fatal("Unable to start:", err)
	}
}
