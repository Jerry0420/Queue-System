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

    var db *sql.DB
    dbLocation := fmt.Sprintf("%s:%d/%s?sslmode=%s", 
        serverConfig.POSTGRES_HOST(), 
        serverConfig.POSTGRES_PORT(), 
        serverConfig.POSTGRES_DB(),
        serverConfig.POSTGRES_SSL(),
    )
    if serverConfig.ENV() == "prod" {
        logical, token, sys := config.NewVaultConnection(
            serverConfig.VAULT_SERVER(), 
            serverConfig.VAULT_WRAPPED_TOKEN_SERVER(),
            serverConfig.VAULT_ROLE_ID(),
            serverConfig.VAULT_CRED_NAME(),
            logger,
        )
        vaultWrapper := config.NewVaultWrapper(
            serverConfig.VAULT_CRED_NAME(),
            logical, 
            token, 
            sys,
            logger,
        )
        dbWrapper := repository.NewDbWrapper(vaultWrapper, dbLocation, logger)
        db = dbWrapper.GetDb()
        go func() {
            err := dbWrapper.ClosdAllDbConns()
            if err != nil {
                logger.ERRORf("db connection close fail %v", err)
            }
        }()
        go func() {
            err := vaultWrapper.RevokeToken()
            if err != nil {
                logger.WARNf("Fail to revoke token. %v", err)
            }
        }()
    } else {
        db = repository.GetDevDb(serverConfig.POSTGRES_DEV_USER(), serverConfig.POSTGRES_DEV_PASSWORD(), dbLocation, logger)
        go func() {
            err := db.Close()
            if err != nil {
                logger.ERRORf("dev db connection close fail %v", err)
            }
        }()
    }

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
        logger.INFOf("Server Start!")
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