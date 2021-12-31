package usecase

import (
	"context"
	"fmt"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
)

type sessionUsecase struct {
	pgDBSessionRepository  pgDB.PgDBSessionRepositoryInterface
	logger                 logging.LoggerTool
}

func NewSessionUsecase(
	pgDBSessionRepository pgDB.PgDBSessionRepositoryInterface,
	logger logging.LoggerTool,
) SessionUseCaseInterface {
	return &sessionUsecase{pgDBSessionRepository, logger}
}

func (su *sessionUsecase) CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error) {
	session, err := su.pgDBSessionRepository.CreateSession(ctx, store)
	return session, err
}

func (su *sessionUsecase) UpdateSessionStatus(ctx context.Context, session *domain.StoreSession, oldStatus string, newStatus string) error {
	err := su.pgDBSessionRepository.UpdateSessionStatus(ctx, nil, session, oldStatus, newStatus)
	session.StoreSessionStatus = newStatus
	return err
}

func (su *sessionUsecase) TopicNameOfUpdateSession(storeId int) string {
	return fmt.Sprintf("updateSession.%d", storeId)
}

func (su *sessionUsecase) GetSessionById(ctx context.Context, sessionId string) (session domain.StoreSession, err error) {
	if sessionId == "" {
		return session, domain.ServerError40106
	}
	session, err = su.pgDBSessionRepository.GetSessionById(ctx, sessionId)
	return session, err
}
