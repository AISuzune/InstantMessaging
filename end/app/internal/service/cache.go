package service

import (
	g "InstantMessaging/app/global"
	"context"
	"time"
)

func SetUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	ctx := context.Background()
	g.Rdb.Set(ctx, key, val, timeTTL)
}
