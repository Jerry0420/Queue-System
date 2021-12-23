package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/jerry0420/queue-system/backend/domain"
)

func (uc *Usecase) CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error) {
	session, err := uc.pgDBRepository.CreateSession(ctx, store)
	return session, err
}

func (uc *Usecase) UpdateSessionStatus(ctx context.Context, session *domain.StoreSession, oldStatus string, newStatus string) error {
	err := uc.pgDBRepository.UpdateSessionStatus(ctx, session, oldStatus, newStatus)
	session.StoreSessionStatus = newStatus
	return err
}

func (uc *Usecase) TopicNameOfUpdateSession(storeId int) string {
	return fmt.Sprintf("updateSession.%d", storeId)
}

func (uc *Usecase) GetSessionAndStoreBySessionId(ctx context.Context, sessionId string) (session domain.StoreSession, store domain.Store, err error) {
	if sessionId == "" {
		return session, store, domain.ServerError40106
	}
	session, store, err = uc.pgDBRepository.GetSessionAndStoreBySessionId(ctx, sessionId)
	store, err = uc.CheckStoreExpirationStatus(store, err)
	switch {
	case store == domain.Store{} && err != nil:
		return session, store, err
	case store != domain.Store{} && errors.Is(err, domain.ServerError40903):
		_ = uc.CloseStore(ctx, store)
		return session, store, err
	}

	return session, store, nil
}
