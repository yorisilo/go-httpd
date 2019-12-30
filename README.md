2019-12-30 02:43:52

# go-httpd
cf. https://ascii.jp/elem/000/001/276/1276572/

httpd(webサーバ) を golang かつ TCPの機能を使って作ってみる

golang で httpd を作成するときは普通は高機能な net/http以下のAPI(アプリケーション層のプロトコルをしゃべるAPI)を使うことが多いが、
どうやってwebサーバが実現されているかを知るために自分でhttpdを作ってみるのが良さそう。ソケット通信を直接扱う低レイヤーなnetのAPI(トランスポート層のTCPなどをしゃべるAPI)を使って作ってみる。

net/http APIもその下の net API を使ってhttp通信を実現している。

## TCPの機能(net.Conn)だけを使って webサーバを作ってみる
