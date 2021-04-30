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


※ここでは、[参考サイト](https://outcrawl.com/go-graphql-realtime-chat)のなかで、`subscription`に関する部分だけ書き出す。  

```graphql
scalar Time

type Message {
  id: String!
  user: String!
  createdAt: Time!
  text: String!
}

type Subscription {
  messagePosted(user: String!): Message!
  userJoined(user: String!): String!
}
```

続いて、`./graph/schema.resolvers.go`を一度削除する。  
その後、`gqlgen generate`を実行する。  
※`./graph/schema.resolvers.go`を削除しない場合、エラーになる。

