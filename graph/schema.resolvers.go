package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SRsawaguchi/graphql-subscription-server/graph/generated"
	"github.com/SRsawaguchi/graphql-subscription-server/graph/model"
	"github.com/segmentio/ksuid"
)

func (r *mutationResolver) PostMessage(ctx context.Context, user string, text string) (*model.Message, error) {
	message := &model.Message{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		User:      user,
		Text:      text,
	}

	// 投稿されたメッセージを保存し、subscribeしている全てのコネクションにブロードキャスト
	r.mutex.Lock()
	r.messages = append(r.messages, message)
	for _, ch := range r.subscribers {
		ch <- message
	}
	r.mutex.Unlock()

	return message, nil
}

func (r *queryResolver) Messages(ctx context.Context) ([]*model.Message, error) {
	return r.messages, nil
}

func (r *subscriptionResolver) MessagePosted(ctx context.Context, user string) (<-chan *model.Message, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.subscribers[user]; ok {
		err := fmt.Errorf("`%s` has already been subscribed.", user)
		log.Print(err.Error())
		return nil, err
	}

	// チャンネルを作成し、リストに登録
	ch := make(chan *model.Message, 1)
	r.subscribers[user] = ch
	log.Printf("`%s` has been subscribed!", user)

	// コネクションが終了したら、このチャンネルを削除する
	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.subscribers, user)
		r.mutex.Unlock()
		log.Printf("`%s` has been unsubscribed.", user)
	}()

	return ch, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
