package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/usecase"
)

type Middleware struct {
	usecase *usecase.Usecase
	logger  logging.LoggerTool
}

func NewMiddleware(router *mux.Router, logger logging.LoggerTool, usecase *usecase.Usecase) *Middleware {
	mw := &Middleware{usecase, logger}
	router.Use(mw.LoggingMiddleware)
	return mw
}

func (mw *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: remove after dev...
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		start := time.Now()

		randomUUID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "requestID", randomUUID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		ctx = context.WithValue(r.Context(), "duration", time.Since(start).Truncate(1*time.Millisecond))
		if errorCode := w.Header().Get("Server-Code"); errorCode != "" {
			ctx = context.WithValue(ctx, "code", errorCode)
		} else {
			ctx = context.WithValue(ctx, "code", strconv.Itoa(200))
		}

		r = r.WithContext(ctx)
		mw.logger.INFOf(r.Context(), "response")
	})
}

func (mw *Middleware) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encryptToken := strings.Split(r.Header.Get("Authorization"), " ")
		if len(encryptToken) == 2 && strings.ToLower(encryptToken[0]) == "bearer" {
			tokenClaims, err := mw.usecase.VerifyToken(r.Context(), encryptToken[1], domain.SignKeyTypes.NORMAL, mw.usecase.GetSignKeyByID)
			if err != nil {
				presenter.JsonResponse(w, nil, err)
				return
			}

			ctx := context.WithValue(r.Context(), domain.SignKeyTypes.NORMAL, tokenClaims)
			mw.logger.INFOf(ctx, "storeID: %d", tokenClaims.StoreID)
			r = r.WithContext(ctx)

		} else {
			presenter.JsonResponse(w, nil, domain.ServerError40102)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// for customers....
func (mw *Middleware) SessionAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get("Authorization")
		if sessionId != "" {
			session, store, err := mw.usecase.GetSessionAndStoreBySessionId(r.Context(), sessionId)
			store, err = mw.usecase.CheckStoreExpirationStatus(store, err)
			switch {
			case store == domain.Store{} && err != nil:
				presenter.JsonResponse(w, nil, err)
				return
			case store != domain.Store{} && errors.Is(err, domain.ServerError40903):
				_ = mw.usecase.CloseStore(r.Context(), store)
				presenter.JsonResponse(w, nil, domain.ServerError40903)
				return
			}

			ctx := context.WithValue(r.Context(), domain.StoreSessionString, session)
			r = r.WithContext(ctx)
		} else {
			presenter.JsonResponse(w, nil, domain.ServerError40106)
			return
		}
		next.ServeHTTP(w, r)
	})
}
