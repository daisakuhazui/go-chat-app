## 概要
```
- go version go1.11 darwin/amd64
- クライアント間でチャットをすることができます
- チャットを行うには外部サービスでの認証が必要です
- 対応している外部サービスはGoogle です(facebook と GitHub にも順次対応予定)
```

## パッケージのインストール
```
$ cd chat
$ make dep
```

## ビルド方法
```
$ cd chat
$ make build
```

## Webサーバー起動
```
$ cd chat
$ make run
```

## テスト実行
```
$ cd chat
$ make test
```