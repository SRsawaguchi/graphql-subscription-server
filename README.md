# graphql-subscription-server

GraphQLのSubscriptionを試してみる。  
ここではGoのServerを実装する。

- [gqlgen](https://gqlgen.com/getting-started/)を使う。

## 準備

gqlgenが推奨するディレクトリ構造で初期化する。

```
go mod init github.com/SRsawaguchi/graphql-subscription-server
go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen init
```

