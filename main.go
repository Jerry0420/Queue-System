package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"
    "fmt"
    "embed"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/jerry0420/queue-system/backend/config"
	"github.com/jerry0420/queue-system/backend/delivery/http"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/middleware"
	"github.com/jerry0420/queue-system/backend/presenter"
	"github.com/jerry0420/queue-system/backend/repository/db"
	"github.com/jerry0420/queue-system/backend/usecase"
	"github.com/jerry0420/queue-system/backend/utils"
)

//go:embed build
var files embed.FS

func main() {
    logger := logging.NewLogger([]string{"method", "url", "code", "sep", "requestID", "duration"}, false)
    
    serverConfig := config.NewConfig(logger)
    vault := config.NewVault(
        serverConfig.VAULT_SERVER(), 
        serverConfig.VAULT_TOKEN(), 
        serverConfig.VAULT_CRED_NAME(), 
        logger,
    )

    leaseId, username, password := vault.GetDbSecret()
    defer vault.RevokeLease(leaseId)
    dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", 
		username, 
		password, 
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
			logger.ERRORf("db connection close fail %v", err)
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

    frontendFiles := utils.GetFrontendFiles(files, "build")
    delivery.NewFrontendDelivery(router, logger, frontendFiles)

    server := &http.Server{
        Addr:         "0.0.0.0:8000",
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: router,
    }

    go func() {
        if err := server.ListenAndServe(); err != nil {
            logger.ERRORf("ListenAndServe http fail %v", err)
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
}