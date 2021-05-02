# graphql-subscription-server

GraphQLのSubscriptionを試してみる。  
ここではGoのServerを実装する。

- [https://outcrawl.com/go-graphql-realtime-chat](https://outcrawl.com/go-graphql-realtime-chat)を参考にする。
- [gqlgen](https://gqlgen.com/getting-started/)を使う。


## 準備

gqlgenが推奨するディレクトリ構造で初期化する。

```
go mod init github.com/SRsawaguchi/graphql-subscription-server
go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen init
```


## schemaの定義
GraphQLのschemaを定義する。  
`./graph/schema.graphqls`を編集する。  


※ここでは、[参考サイト](https://outcrawl.com/go-graphql-realtime-chat)のなかで、必要な部分だけ書き出す。  

```graphql
scalar Time

type Message {
  id: String!
  user: String!
  createdAt: Time!
  text: String!
}

type Mutation {
  postMessage(user: String!, text: String!): Message
}

type Query {
  messages: [Message!]!
}

type Subscription {
  messagePosted(user: String!): Message!
}
```

続いて、`./graph/schema.resolvers.go`を一度削除する。  
その後、`gqlgen generate`を実行する。  
※`./graph/schema.resolvers.go`を削除しない場合、エラーになる。
※この後、スキーマを書き変えた場合は、再び`gqlgen generate`すれば`schema.resolvers.go`が自動生成される。  


## resolver.goの実装
`./graph/resolver.go`には`Resolver`構造体の定義だけがある。  
ここはそのアプリに特有の依存関係などをResolver構造体に自由に追加できる。  
つまり、ここでDIするということ。 

```go
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
```

今回は`NewResolver()`という関数を作成。  
これを`server.go`から利用する。  

```go
// 省略
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver()}))
// 省略
```

## schema.resolvers.goの実装
`gqlgen generate`を行うと、`./graph/schema.resolvers.go`が自動で生成される。  
これは各種GraphQLのQueryやMutation、Subscriptionのエントリポイントのひな形が定義されている。  
このひな形のなかに、各ビジネスロジックを書いていく。  

例えば、サブスクリプションである`messagePosted`はこのような感じになる。  

```go
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
```

## サーバを起動
以下のコマンドでサーバを起動する。  

```
go run server.go
```

## クエリ
以下のようなクエリを実行する。  
※GraphQL playgroundを利用するとよい。タブを2つ使って、片方はsubscription、片方はmutation(postMessage)を実行するとよい。2つのブラウザを使ってやると、メッセージがブロードキャストされていることを確認できる。

### messagePosted(サブスクリプション)
```graphql
subscription($user: String!) {
  messagePosted(user: $user) {
    id
    user
    text
    createdAt
  }
}
```

valiables(例)
```
{
  "user": "tanaka"
}
```

### postMessage(クエリ)
```graphql
mutation($user: String!, $text: String!) {
  postMessage(user: $user, text: $text) {
    id
    user
    text
    createdAt
  }
}
```
valiables(例)
```
{
  "user": "tanaka",
  "text": "Hi there!!"
}
```
