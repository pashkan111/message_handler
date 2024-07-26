package main

import (
	"context"
	"fmt"
	"messange_handler/src/api"
	"messange_handler/src/dependencies/logger"
	"messange_handler/src/dependencies/pg"
	"messange_handler/src/queue"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.GetLogger()
	postgres_pool := pg.GetPostgresPool(ctx, log)

	queue.InitProducer(log)

	message_consumer := queue.MessageConsumer{Log: log}
	message_consumer.Run(ctx, postgres_pool, log)

	router := mux.NewRouter()
	api.InitMessageRoutes(router, postgres_pool, log)

	web_server_port := os.Getenv("WEB_SERVER_PORT")
	port, err := strconv.Atoi(web_server_port)
	if err != nil {
		log.Fatal("Error with parsing WEB_SERVER_PORT", err)
	}

	log.Infof("Server is running on port %d", port)
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("127.0.0.1:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
