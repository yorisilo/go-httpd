# golang で httpd を作ってみた
- TCP の機能(`net.Conn`)だけを使って HTTP による通信を実現してみる

TCP ソケットを使って echo サーバーを作るのではなく、 HTTP サーバーを作ることを目的とする。
