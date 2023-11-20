package simplecache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/singleflight"
)

type Client interface {
	Get(context.Context, string) ([]byte, bool, error)
	Set(context.Context, string, []byte, time.Duration) error
}

type Cache[T any] struct {
	client     Client
	expiration time.Duration
	sfg        *singleflight.Group
}

func NewCache[T any](client Client, expiration time.Duration) *Cache[T] {
	return &Cache[T]{
		client:     client,
		expiration: expiration,
		sfg:        &singleflight.Group{},
	}
}

func (c *Cache[T]) GetOrSet(
	ctx context.Context,
	key string,
	callback func(context.Context) (T, error),
) (T, error) {
	// singleflightでリクエストをまとめる
	a, err, _ := c.sfg.Do(key, func() (any, error) {
		// キャッシュから取得
		bytes, exist, err := c.client.Get(ctx, key)
		if err != nil {
			log.Println(err)
		}
		if exist {
			return bytes, nil
		}
		// キャッシュがなければcallbackを実行
		t, err := callback(ctx)
		if err != nil {
			return nil, err
		}
		bytes, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		// キャッシュに保存
		err = c.client.Set(ctx, key, bytes, c.expiration)
		if err != nil {
			log.Println(err)
		}
		return bytes, nil
	})
	var t T
	if err != nil {
		return t, err
	}
	bytes, ok := a.([]byte)
	if !ok {
		// 実装上、起きることはないはず
		return t, fmt.Errorf("failed to get from cache: invalid type %T", a)
	}
	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}
