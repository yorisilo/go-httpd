# network についての発表
golang でアプリケーション層のライブラリは使わずに、 トランスポート層のみを扱うライブラリを使って httpd を http0.9 -> http1.0 -> http1.1 と徐々に進化させて作成する予定だった。
その過程で色々と疑問に思ったことを調査したりしたのでその話をしたい。

## 三部作
- httpd を作成中に、 `conn.Read` せずに `conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\nHello World\n"))` したら、 curl で `curl: (56) Recv failure: Connection reset by peer` 怒られてしまい、その理由を知りたくなった。結果 RFC 7230 を読む羽目になった。
- ネットワークのこと
- httpd 作成再び。 では、http の仕様がなんとなく把握できたところで httpd を作っていく。
http0.9 -> http1.0 -> http1.1 の仕様に従ってだんだんと httpd を進化させていく
