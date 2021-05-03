package graph

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/SRsawaguchi/graphql-subscription-server/graph/model"
	"github.com/go-redis/redis/v8"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	redisClient *redis.Client
	redisPubSub *redis.PubSub
	subscribers map[string]chan<- *model.Message
	mutex       sync.Mutex
}

const redisPostMessagesSubscription = "messages"
const redisKeyMessages = "messages"

func NewResolver(ctx context.Context) *Resolver {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// messagesチャンネルを購読
	pubsub := redisClient.Subscribe(ctx, redisPostMessagesSubscription)

	resolver := &Resolver{
		redisClient: redisClient,
		redisPubSub: pubsub,
		subscribers: map[string]chan<- *model.Message{},
		mutex:       sync.Mutex{},
	}

	// messagesにpublishされたデータを取得した場合の処理
	// ゴルーチンを使って非同期で行う
	go func() {
		pubsubCh := pubsub.Channel()

		// メッセージの受信（consume）
		for msg := range pubsubCh {
			// 受信したmessageはJSON形式なので、これをmodel.Message構造体に変換
			message := &model.Message{}
			err := json.Unmarshal([]byte(msg.Payload), message)
			if err != nil {
				log.Printf(err.Error())
				continue
			}

			// 購読しているクライアントにRedisから受け取ったMessageをブロードキャスト
			resolver.mutex.Lock()
			for _, ch := range resolver.subscribers {
				ch <- message
			}
			resolver.mutex.Unlock()
		}
	}()

	return resolver
}
