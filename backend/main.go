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
	"google.golang.org/grpc"

	"github.com/jerry0420/queue-system/backend/broker"
	"github.com/jerry0420/queue-system/backend/config"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/middleware"
	"github.com/jerry0420/queue-system/backend/logging"

	"github.com/jerry0420/queue-system/backend/repository/grpcServices"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
	"github.com/jerry0420/queue-system/backend/usecase"
)

func main() {
	logger := logging.NewLogger([]string{"method", "url", "code", "sep", "requestID", "duration"}, false)

	var db *sql.DB
	dbLocation := config.ServerConfig.POSTGRES_LOCATION()
	var grpcConn *grpc.ClientConn
	var grpcClient grpcServices.GrpcServiceClient

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
		grpcConn, grpcClient = grpcServices.GetGrpcConn(logger, config.ServerConfig.GRPC_HOST(), config.ServerConfig.CA_CRT())
		defer func() {
			revokeTokenErr := vaultWrapper.RevokeToken()
			if revokeTokenErr != nil {
				logger.WARNf("Fail to revoke token. %v", revokeTokenErr)
			}
		}()

	} else {
		db = pgDB.GetDevDb(config.ServerConfig.POSTGRES_DEV_USER(), config.ServerConfig.POSTGRES_DEV_PASSWORD(), dbLocation, logger)
		grpcConn, grpcClient = grpcServices.GetDevGrpcConn(logger, config.ServerConfig.GRPC_HOST())
	}

	defer func() {
		err := db.Close()
		if err != nil {
			logger.ERRORf("db connection close fail %v", err)
		}
		err = grpcConn.Close()
		if err != nil {
			logger.ERRORf("grpc connection close fail %v", err)
		}
	}()

	router := mux.NewRouter()

	pgDBTx := pgDB.NewPgDBTx(db, logger)
	pgDBStoreRepository := pgDB.NewPgDBStoreRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	pgDBSignkeyRepository := pgDB.NewPgDBSignKeyRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	pgDBSessionRepository := pgDB.NewPgDBSessionRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	pgDBQueueRepository := pgDB.NewPgDBQueueRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	pgDBCustomerRepository := pgDB.NewPgDBCustomerRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())

	grpcServicesRepository := grpcServices.NewGrpcServicesRepository(grpcClient, logger, config.ServerConfig.CONTEXT_TIMEOUT()*4)

	storeUsecase := usecase.NewStoreUsecase(
		pgDBTx,
		pgDBStoreRepository,
		pgDBSignkeyRepository,
		logger,
		usecase.StoreUsecaseConfig{
			Domain: config.ServerConfig.DOMAIN(),
		},
	)
	sessionUsecase := usecase.NewSessionUsecase(pgDBSessionRepository, logger)
	customerUsecase := usecase.NewCustomerUsecase(pgDBCustomerRepository, logger)
	integrationUsecase := usecase.NewIntegrationUsecase(
		pgDBTx,
		pgDBStoreRepository,
		pgDBSessionRepository,
		pgDBCustomerRepository,
		pgDBQueueRepository,
		grpcServicesRepository,
		storeUsecase,
		logger,
		usecase.IntegrationUsecaseConfig{
			StoreDuration: config.ServerConfig.STOREDURATION(),
			TokenDuration: config.ServerConfig.TOKENDURATION(),
			PasswordTokenDuration: config.ServerConfig.PASSWORDTOKENDURATION(),
			GrpcReplicaCount: config.ServerConfig.GRPCREPLICACOUNT(),
		},
	)

	broker := broker.NewBroker(logger)
	defer broker.CloseAll()

	mw := middleware.NewMiddleware(router, logger, integrationUsecase, sessionUsecase)

	httpAPI.NewHttpAPIDelivery(
		router,
		logger,
		mw,
		customerUsecase,
		sessionUsecase,
		storeUsecase,
		integrationUsecase,
		broker,
		httpAPI.HttpAPIDeliveryConfig{
			StoreDuration:         config.ServerConfig.STOREDURATION(),
			TokenDuration:         config.ServerConfig.TOKENDURATION(),
			PasswordTokenDuration: config.ServerConfig.PASSWORDTOKENDURATION(),
			Domain:                config.ServerConfig.DOMAIN(),
		},
	)

	server := &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Hour * 24,
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
