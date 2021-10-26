package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/config"
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
	db *sql.DB
	nextDb *sql.DB
	leaseId string
	newxtLeaseId string
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

func (dbw *DbWrapper) checkAndRenewDb() {
	
}

func (dbw *DbWrapper) GetDb() *sql.DB {
	go dbw.checkAndRenewDb()
	return dbw.db
}

func (dbw *DbWrapper) ClosdAllDbConns() error {
	// 順便把 lease revoke
	return nil
}