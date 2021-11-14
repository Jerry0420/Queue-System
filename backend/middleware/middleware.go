package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/presenter"
)

type Middleware struct {
	storeUsecase domain.StoreUsecaseInterface
	logger       logging.LoggerTool
}

func NewMiddleware(router *mux.Router, logger logging.LoggerTool, storeUsecase domain.StoreUsecaseInterface) *Middleware {
	mw := &Middleware{storeUsecase, logger}
	router.Use(mw.LoggingMiddleware)
	return mw
}

func (mw *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		randomUUID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "requestID", randomUUID)

		r = r.WithContext(ctx)
		responseWrapper := &presenter.ResponseWrapper{ResponseWriter: w, Buffer: &bytes.Buffer{}}
		next.ServeHTTP(responseWrapper, r)

		var wrappedResponse map[string]interface{}
		json.Unmarshal(responseWrapper.Buffer.Bytes(), &wrappedResponse)
		io.Copy(w, responseWrapper.Buffer)

		ctx = context.WithValue(r.Context(), "duration", time.Since(start).Truncate(1*time.Millisecond))

		if errorCode, ok := wrappedResponse["error_code"]; ok {
			// api routes will go here.
			ctx = context.WithValue(ctx, "code", errorCode)
		} else {
			ctx = context.WithValue(ctx, "code", 200)
		}

		r = r.WithContext(ctx)
		mw.logger.INFOf(r.Context(), "response")
	})
}

func (mw *Middleware) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encryptToken := r.Header.Get("Authorization")
		if len(encryptToken) > 0 {
			tokenClaims, err := mw.storeUsecase.VerifyToken(r.Context(), encryptToken, mw.storeUsecase.GetSignKeyByID)
			switch {
			case tokenClaims != domain.TokenClaims{} && errors.Is(err, domain.ServerError40104):
				store := domain.Store{ID: tokenClaims.StoreID, Email: tokenClaims.Email, Name: tokenClaims.Name}
				err = mw.storeUsecase.Close(r.Context(), store)
				presenter.JsonResponse(w, nil, domain.ServerError40903)
				return
			case err != nil:
				presenter.JsonResponse(w, nil, err)
				return
			}

			ctx := context.WithValue(r.Context(), "token", tokenClaims)
			mw.logger.INFOf(ctx, "storeID: %d", tokenClaims.StoreID)
			r = r.WithContext(ctx)

		} else {
			presenter.JsonResponse(w, nil, domain.ServerError40102)
			return
		}
		next.ServeHTTP(w, r)
	})
}
