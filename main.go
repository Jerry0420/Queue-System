package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"

	"github.com/jerry0420/queue-system/config"
	"github.com/jerry0420/queue-system/logging"
	"github.com/jerry0420/queue-system/domain"
    "github.com/jerry0420/queue-system/presenter"
    "github.com/jerry0420/queue-system/repository/db"
    "github.com/jerry0420/queue-system/usecase"
    "github.com/jerry0420/queue-system/delivery/http"
    "github.com/jerry0420/queue-system/middleware"
)

func main() {
    logger := logging.NewLogger([]string{"method", "url", "code", "sep", "requestID", "duration"}, false)
    serverConfig := config.NewConfig()
    dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", 
        serverConfig.POSTGRES_USER(), 
        serverConfig.POSTGRES_PASSWORD(), 
        serverConfig.POSTGRES_HOST(), 
        serverConfig.POSTGRES_PORT(), 
        serverConfig.POSTGRES_DB(),
        serverConfig.POSTGRES_SSL(),
    )
    
    db, err := sql.Open("postgres", dbConnectionString)
    if err != nil {
        logger.FATALf("db connection fail %v", err)
    }
    
    err = db.Ping()
    if err != nil {
        logger.FATALf("db ping fail %v", err)
    }
    
    defer func() {
		err := db.Close()
		if err != nil {
			logger.FATALf("db connection close fail %v", err)
		}
	}()

    router := mux.NewRouter()

    storeReposotory := repository.NewStoreRepository(db, logger)
    queueReposotory := repository.NewQueueRepository(db, logger)
    customerReposotory := repository.NewCustomerRepository(db, logger)

    storeUsecase := usecase.NewStoreUsecase(storeReposotory, logger)
    queueUsecase := usecase.NewQueueUsecase(queueReposotory, logger)
    customerUsecase := usecase.NewCustomerUsecase(customerReposotory, logger)

    middleware.NewMiddleware(router, logger)

    delivery.NewStoreDelivery(router, logger, storeUsecase)
    delivery.NewQueueDelivery(router, logger, queueUsecase)
    delivery.NewCustomerDelivery(router, logger, customerUsecase)

    router.HandleFunc("/hello", func (w http.ResponseWriter, r *http.Request)  {
        presenter.JsonResponseOK(w, map[string]string{"hello": "world"})
    })

    // final route, for unsupported route!.
    router.HandleFunc("/{rest_of_router}", func (w http.ResponseWriter, r *http.Request)  {
        presenter.JsonResponse(w, nil, domain.ServerError40401)
    })

    server := &http.Server{
        Addr:         "0.0.0.0:8000",
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: router,
    }

    go func() {
        if err := server.ListenAndServe(); err != nil {
            logger.FATALf("ListenAndServe http fail %v", err)
        }
    }()
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)

    // Block until receive signal...
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
    defer cancel()
    server.Shutdown(ctx)
    logger.INFOf("shutting down")
    os.Exit(0)
}