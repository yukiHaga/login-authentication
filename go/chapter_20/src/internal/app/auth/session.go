package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yukiHaga/web_server/src/internal/app/model"
)

type SessionId string

type Session struct {
	Store *redis.Client
}

var ctx = context.Background()

func NewSession() *Session {
	store := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0, // 確かデフォルトのdbを使うみたいな意味だった
	})

	return &Session{Store: store}
}

func (session *Session) Save(userId model.UserId) (SessionId, error) {
	sessionId := uuid.NewString()
	// 0は確か有効期限とか無かった気がする
	return SessionId(sessionId), session.Store.Set(ctx, sessionId, int64(userId), 0).Err()
}

func (session *Session) Load(sessionId string) (model.UserId, error) {
	id, err := session.Store.Get(ctx, sessionId).Int64()
	if err != nil {
		return 0, err
	}

	return model.UserId(id), nil
}
