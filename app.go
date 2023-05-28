package main

import (
	"context"
	"fmt"
	"main/common/config"
	"main/middlewares"
	"main/routes"
	"main/services"
	"main/utils/logging"
	"main/utils/mongodb"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger()

	mongoDBClient, err := mongodb.NewClient(context.Background(),
		cfg.MongoDB.Host, cfg.MongoDB.Port, cfg.MongoDB.Username,
		cfg.MongoDB.Password, cfg.MongoDB.Database, cfg.MongoDB.AuthDB,
	)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Passwod,
		DB:       cfg.Redis.DB,
	})
	rp := rdb.Ping(context.Background())
	if rp.Err() != nil {
		panic(err)
	}

	logger.Info("Create services")
	services := services.NewServices(mongoDBClient, logger)

	logger.Info("Create middlewares")
	middlewares := middlewares.NewMiddlewares(rdb, logger)

	logger.Info("Create handler")
	router := routes.NewRouter(services, rdb, middlewares, logger)

	logger.Info("Create router and register routes")
	router.Register()

	start(router.Router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()

	var listener net.Listener
	var listenerErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Create socket")

		socketPath := path.Join(appDir, "app.sock")

		logger.Debugf("Socket path: %s", socketPath)

		listener, listenerErr = net.Listen("unix", socketPath)

		logger.Infof("Start server on unix socket: %s", socketPath)
	} else {
		logger.Infof("Listen %s", cfg.Listen.Type)

		listener, listenerErr = net.Listen(cfg.Listen.Type, fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))

		logger.Infof("Start server on %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}

	if listenerErr != nil {
		logger.Fatal(listenerErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
