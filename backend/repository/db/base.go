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

type DbWrapper struct {
	sync.Mutex
	db *sql.DB
	leaseRevocableChan chan bool
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
	var credExpireChan chan bool
	var dbConnectionString string
	var db *sql.DB
	var err error
	leaseRevocableChan := make(chan bool, 1)

	for {
		username, password, credExpireChan = dbw.vaultWrapper.GetDbCred(leaseRevocableChan)
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

		if dbw.db == nil {
			dbw.db = db
			dbw.leaseRevocableChan = leaseRevocableChan
		} else {
			dbw.Lock()
			// close old db connection and set to new db connnection.
			err = dbw.db.Close()
			if err != nil {
				dbw.logger.ERRORf("db connection close fail %v", err)
			}
			dbw.db = db
			dbw.Unlock()
			dbw.leaseRevocableChan <- true
			dbw.leaseRevocableChan = leaseRevocableChan
		}
		<- credExpireChan
	}
}

func (dbw *DbWrapper) GetDb() *sql.DB {
	go dbw.CheckAndRenewDb()
	for {
		if dbw.db != nil {
			return dbw.db
		}
	}
}
