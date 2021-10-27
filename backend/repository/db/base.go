package repository

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/jerry0420/queue-system/backend/config"
	"github.com/jerry0420/queue-system/backend/logging"
	_ "github.com/lib/pq"
)

func GetDevDb(username string, password string, dbLocation string, logger logging.LoggerTool) *sql.DB {
	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s", 
		username, 
		password, 
		dbLocation,
    )
    
    db, err := sql.Open("postgres", dbConnectionString)
    if err != nil {
        logger.FATALf("db connection fail %v", err)
    }
    
    err = db.Ping()
    if err != nil {
        logger.FATALf("db ping fail %v", err)
    }
	return db
}

type dbRotateInfo struct {
	sync.Mutex
	db *sql.DB
	leaseId string
	credExpireChan chan bool
}

type DbWrapper struct {
	dbRI *dbRotateInfo
	dbLocation string
	vaultWrapper *config.VaultWrapper
	logger logging.LoggerTool
}

func NewDbWrapper(vaultWrapper *config.VaultWrapper, dbLocation string, logger logging.LoggerTool) *DbWrapper {
	return &DbWrapper{
		dbLocation: dbLocation,
		vaultWrapper: vaultWrapper,
		logger: logger,
	}
}

func (dbw *DbWrapper) CheckAndRenewDb() {
	var username string
	var password string
	var leaseId string
	var credExpireChan chan bool
	var dbConnectionString string
	var db *sql.DB
	var err error

	for {
		username, password, leaseId, credExpireChan = dbw.vaultWrapper.GetDbCred()
		dbConnectionString = fmt.Sprintf("postgres://%s:%s@%s", 
			username, 
			password, 
			dbw.dbLocation,
		)
		db, err = sql.Open("postgres", dbConnectionString)
		if err != nil {
			dbw.logger.FATALf("db connection fail %v", err)
		}
		
		err := db.Ping()
		if err != nil {
			dbw.logger.FATALf("db ping fail %v", err)
		}

		if dbw.dbRI == nil {
			dbw.dbRI = &dbRotateInfo{db: db, leaseId: leaseId, credExpireChan: credExpireChan}
		} else {
			dbw.dbRI.Lock()
			err = dbw.dbRI.db.Close()
			if err != nil {
				dbw.logger.ERRORf("db connection close fail %v", err)
			}
			dbw.dbRI.db = db
			dbw.dbRI.Unlock()

			dbw.vaultWrapper.RevokeLease(dbw.dbRI.leaseId)
			dbw.dbRI.leaseId = leaseId
			dbw.dbRI.credExpireChan = credExpireChan
		}
		<- dbw.dbRI.credExpireChan
	}
}

func (dbw *DbWrapper) GetDb() *sql.DB {
	go dbw.CheckAndRenewDb()
	for {
		if dbw.dbRI != nil {
			return dbw.dbRI.db
		}
	}
}
