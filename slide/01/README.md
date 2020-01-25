slidenumbers: true
footer: intro
autoscale: true
build-lists: true
theme: Simple, 1

# golang で httpd を作成してみた話 その１
**いかにして私は RFC7230 を読むハメになったのか**

---

# 三行で説明すると

- httpd を作成中に、request を `conn.Read` せずに `conn.Write` した
- curl で `curl: (56) Recv failure: Connection reset by peer` と怒られてしまい、その理由を知りたくなった。
- RFC 7230 を読む羽目になった。

---

http サーバーを

- アプリケーションレイヤを扱うパッケージで作るのではなく、
- トランスポートレイヤのパッケージのみで作ってみたい

という動機があった。

---

- http 接続がどう行われてるか図で説明する

---

ではそれを実装していく

---

- HTTP レイヤのパッケージを使った実装
- TCP ソケットを直接使った実装
とを比べる

---

TCP ソケットを使った実装で read(2) しなかったら `connection reset by peer` と怒られてしまった。調査してみる

- client の request を読まなかった場合
- client の request を読みはしたが全て読まなかった場合
- ciletn の request を全て読んだ場合

---

それぞれのパケットキャプチャ結果を比べてみる

---

- tcpdump の小ネタを挟む。 ACKフラグが `.` になってて混乱することとか。
- tshark はパケットキャプチャ結果が見やすくていいぞとか

---

よくわからん。
ツイッターで質問
RFC 7230 にはこう書いてますよ

---

RFC 7230 を読む

---

なるほど。納得

---

おわり。
