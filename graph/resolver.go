package graph

import (
	"sync"

	"github.com/SRsawaguchi/graphql-subscription-server/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	subscribers map[string]chan<- *model.Message
	messages    []*model.Message
	mutex       sync.Mutex
}

func NewResolver() *Resolver {
	return &Resolver{
		subscribers: map[string]chan<- *model.Message{},
		mutex:       sync.Mutex{},
	}
}
