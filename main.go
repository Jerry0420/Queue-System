package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/jerry0420/queue-system/backend/config"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/middleware"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
	"github.com/jerry0420/queue-system/backend/usecase"
)

func main() {
	logger := logging.NewLogger([]string{"method", "url", "code", "sep", "requestID", "duration"}, false)

	var db *sql.DB
	dbLocation := config.ServerConfig.POSTGRES_LOCATION()

	if config.ServerConfig.ENV() == config.EnvStatus.PROD {
		vaultConnectionConfig := config.ServerConfig.VAULT_CONNECTION_CONFIG()
		logical, token, sys := config.NewVaultConnection(logger, &vaultConnectionConfig)
		vaultWrapper := config.NewVaultWrapper(
			config.ServerConfig.VAULT_CRED_NAME(),
			logical,
			token,
			sys,
			logger,
		)
		dbWrapper := pgDB.NewDbWrapper(vaultWrapper, dbLocation, logger)
		db = dbWrapper.GetDb()

		defer func() {
			dbCloseErr := db.Close()
			if dbCloseErr != nil {
				logger.ERRORf("db connection close fail %v", dbCloseErr)
			}
			revokeTokenErr := vaultWrapper.RevokeToken()
			if revokeTokenErr != nil {
				logger.WARNf("Fail to revoke token. %v", revokeTokenErr)
			}
		}()
	} else {
		db = pgDB.GetDevDb(config.ServerConfig.POSTGRES_DEV_USER(), config.ServerConfig.POSTGRES_DEV_PASSWORD(), dbLocation, logger)
		defer func() {
			err := db.Close()
			if err != nil {
				logger.ERRORf("dev db connection close fail %v", err)
			}
		}()
	}

	router := mux.NewRouter()

	signKeyReposotory := pgDB.NewSignKeyRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	storeReposotory := pgDB.NewStoreRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	queueReposotory := pgDB.NewQueueRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	customerReposotory := pgDB.NewCustomerRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())

	storeUsecase := usecase.NewStoreUsecase( 
		storeReposotory, 
		signKeyReposotory,
		queueReposotory,
		logger,
		usecase.StoreUsecaseConfig{
			Domain: config.ServerConfig.DOMAIN(),
			StoreDuration: config.ServerConfig.STOREDURATION(),
		},
	)
	queueUsecase := usecase.NewQueueUsecase(queueReposotory, logger)
	customerUsecase := usecase.NewCustomerUsecase(customerReposotory, logger)

	mw := middleware.NewMiddleware(router, logger, storeUsecase)
	
	httpAPI.NewStoreDelivery(
		router, 
		logger, 
		mw, 
		storeUsecase, 
		httpAPI.StoreDeliveryConfig{
			StoreDuration: config.ServerConfig.STOREDURATION(), 
			TokenDuration: config.ServerConfig.TOKENDURATION(), 
			PasswordTokenDuration: config.ServerConfig.PASSWORDTOKENDURATION(),
		},
	)
	httpAPI.NewQueueDelivery(router, logger, queueUsecase)
	httpAPI.NewCustomerDelivery(router, logger, customerUsecase)
	httpAPI.NewBaseDelivery(router, logger)

	server := &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	server.Shutdown(ctx)
	logger.INFOf("shutting down")
	os.Exit(0)
}
