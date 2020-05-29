2020-02-08 03:51:31

# TCP/IP の息吹を感じよ
- keywords: netstat, lsof, ifconfig, tcpdump(tshark)

## この発表について
### 書いた人
- ネットワークに詳しくないアプリケーションエンジニア
  - ネットワークキャプチャとかあんまり知らんなあって感じ

### 対象者
- ネットワークに詳しくないアプリケーションエンジニア

### ゴール
- ネットワーク系のコマンドを通じて TCP/IP ネットワーク上を流れるデータの中身の雰囲気がつかめるようになる

## ３行(+α)で説明する
- ネットワーク上を行ったり来たりしてるデータってどんなものか知りたい
- ローカルで http サーバーを立ち上げて、リクエスト送り、レスポンスが返ってくるその経路上で何が行われているか調べてみよう。
- まずは、 http サーバーを立ち上げる
- netstat でネットワークの状態を確かめて、 LISTEN している tcp ソケットを調べる => http サーバーに割り当てられた ip と port が LISTEN してるか確認する
- lsof でそのポートで使用しているプロセスを調べる => ポートに割り当てられているプロセスが http サーバーになっているか確認する
- ifconfig で NIC に割り当てられている IP や MACアドレスを調べる => http サーバーの IP アドレスの NIC を確認する
- tcpdump(tshark) で NIC を指定し、そのネットワーク上を流れているデータを確認し TCP/IP の息吹を感じる

後述するが、BSD 系 OS の Mac のコマンドを使うので Linux(gnu) とはコマンドのオプションなどが違ったり、 Linux にとっては一部非推奨コマンドを使用している。ただ、ネットワークのことについては Linux に対しても使える知識なので安心されたし。

## httpサーバー を建ててみる
http で get リクエストを行うと HelloWorld を返すだけの http サーバー

``` shell
go get -u github.com/shurcooL/goexec
goexec 'http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){io.WriteString(w, "HellWorld\n")}))'
```

## netstat
(トランスポート(tcpとかudp)層の)ネットワークの状態を確かめるコマンド

- ネットワーク接続
- ルーティング
- NIC の状態
などが把握できる

### usage

``` shell
# -a: 全てのアクティブなソケットを表示する
$ netstat -a | grep LISTEN
tcp46      0      0  *.http-alt             *.*                    LISTEN
tcp4       0      0  *.55048                *.*                    LISTEN
tcp4       0      0  localhost.17603        *.*                    LISTEN
tcp4       0      0  localhost.17600        *.*                    LISTEN
tcp6       0      0  *.17500                *.*                    LISTEN
tcp4       0      0  *.17500                *.*                    LISTEN
...

# -n: ホストやユーザーの名前解決を行わず数字のまま出力する
$ netstat -an | grep LISTEN # listen してる全てのアクティブなソケットを表示する
Active Internet connections (including servers)
Proto Recv-Q Send-Q  Local Address          Foreign Address        (state)
tcp46      0      0  *.8080                 *.*                    LISTEN
tcp4       0      0  *.55048                *.*                    LISTEN
tcp4       0      0  127.0.0.1.17603        *.*                    LISTEN
tcp4       0      0  127.0.0.1.17600        *.*                    LISTEN
tcp6       0      0  *.17500                *.*                    LISTEN
tcp4       0      0  *.17500                *.*                    LISTEN
...
```

ちなみに 8080 ポートは http-alt サービスが対応していることは `/etc/services` で確認できる

<details>

``` shell
$ cat /etc/services | grep 8080
http-alt        8080/udp     # HTTP Alternate (see port 80)
http-alt        8080/tcp     # HTTP Alternate (see port 80)
```

</details>

- Proto: プロトコル。tcp とか udp とか。4,6は ipv4 や ipv6 を表す
- Local Address: 接続元のIPとポート。bsd版ではなぜか IP と Port の区切りが `:` でなく `.` である。
- Foreign Address: 接続先のIPとポート
- state: ソケットの状態。 ESTABLISHED(接続が確立されて、通信が行われている) LISTEN(待受状態) TIME_WAIT(接続終了待ち、しばらくするとソケットは閉じられる) などがある

mac(bsd) の netstat だと `netstat -p`(gnu版) ができないので lsof でプロセスを調べる必要がある
bsd 版と gnu 版でだいぶオプション内容が違うので注意

cf.
- [netstatコマンドとは？：ネットワーク管理の基本Tips - ＠IT](https://www.atmarkit.co.jp/ait/articles/1412/10/news003.html)
- [netstat(1) FreeBSDドキュメントJMan](https://kaworu.jpn.org/doc/FreeBSD/jman/man1/netstat.1.php)

## lsof
PortやPID、プロセス名がオープンしているファイルディスクリプタを調べるコマンド

`netstat -an | grep LISTEN` して、 LISTEN してる PORT を調べて、その PORT が使用しているプロセスを調べるときに使ったりする。

### usage

``` shell
# -i:8080:  port 8080 のソケットを対象にする
# -P: ポート名の代わりにポート番号を表示する
# -n: ホスト名の代わりにIPアドレスを表示する
$ lsof -i:8080 -nP # port 8080 で使用しているファイルディスクリプタを調べる
COMMAND   PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
gen     20152 yorisilo    3u  IPv6 0x1bf711a4927fb03b      0t0  TCP *:8080 (LISTEN)

# -p 17018: PID が 17018 を対象にする
$ lsof -p 17018 # PID 17018 がオープンしているファイルディスクリプタを調べる
...
```

cf.
- [開いているファイルのプロセスを特定（lsofコマンド） - Qiita](https://qiita.com/yusabana/items/fd03ee4c90a0d1e0a8c6)
- [【 lsof 】コマンド――オープンしているファイルを一覧表示する：Linux基本コマンドTips（298） - ＠IT](https://www.atmarkit.co.jp/ait/articles/1904/18/news033.html)
- [lsof bsd](https://www.freebsd.org/cgi/man.cgi?query=lsof&manpath=FreeBSD+9.0-RELEASE+and+Ports&format=html)

## ifconfig
ネットワーク・インターフェース(NIC)のネットワーク状態を確認したり、設定を行うコマンド
- NIC に紐付いている IP, MAC アドレスなどを確認できる

### usage

``` shell
# -a: 全てのインタフェースごとのネットワークの状態を表示する
$ ifconfig -a
# : en0 というインターフェースのネットワークの状態を表示する
$ ifconfig en0
en0: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
        options=400<CHANNEL_IO>
        ether AA:BB:CC:DD:EE:FF
        inet6 fe80::ef:a6b3:93b6:fe4a%en0 prefixlen 64 secured scopeid 0x4
        inet 192.168.3.10 netmask 0xffffff00 broadcast 192.168.3.255
        inet6 2400:2410:9340:c200:1c88:1030:a9cd:7a0 prefixlen 64 autoconf secured
        inet6 2400:2410:9340:c200:893f:a7b2:c9a1:fc79 prefixlen 64 autoconf temporary
        nd6 options=201<PERFORMNUD,DAD>
        media: autoselect
        status: active
```

- ether: MAC アドレス
- inet: ip(v4) アドレス
- inet6: ip(v6) アドレス

cf.
- [ifconfig ～（IP）ネットワーク環境の確認／設定を行う：ネットワークコマンドの使い方 - ＠IT](https://www.atmarkit.co.jp/ait/articles/0109/29/news004.html)
- [ターミナルからプライベートIPアドレスとMACアドレス、ルーティングテーブルを確認する - bambinya's blog](http://bambinya.hateblo.jp/entry/2015/04/04/234428)
- [ifconfigの出力結果に書いてあること - Qiita](https://qiita.com/TD3P/items/aff8db72530c6baa11b2)
- [ifconfig(8) FreeBSDドキュメントJMan](https://kaworu.jpn.org/doc/FreeBSD/jman/man8/ifconfig.8.php)

### ネットワーク・インターフェース(NIC)
- ネットワーク接続に必要なインターフェース
  - AWS上でENI(Elastic Network Insterface)に値するもの。
- 物理的なハードウェアでは、NIC(ネットワークインターフェースカード)というカード型の拡張装置を用いる。
  - LAN ケーブルを差し込むハードウェア。アレが物理版の NIC。 LAN カードとかネットワークカードと呼ばれたりもする。
- IPアドレスはNICに設定される
- IPアドレスは「ホスト」に対してではなく、NIC に対して割り当てられる。
- 通常は1NIC-1IPアドレスになるが、1つのNICに複数のIPアドレスを割り当てることが可能だったり、複数のNICを備える機器（ルータなど）が存在する。

ちなみに mac で `ifconfig -a` として出てくる NIC の `awdl0` と `llw0` などは同じ MAC アドレスを持っているので、ハードウェア的に同じものを指していると考えられる。

cf.
- [IPアドレスの基礎知識 - Qiita](https://qiita.com/mogulla3/items/efb4c9328d82d24d98e6#1-3-ip%E3%82%A2%E3%83%89%E3%83%AC%E3%82%B9%E3%81%AFnic%E3%81%AB%E8%A8%AD%E5%AE%9A%E3%81%95%E3%82%8C%E3%82%8B)

#### NIC の仕事内容
- NIC に MAC アドレスが紐付いている
- この MAC アドレスで NIC を識別している
- NIC がイーサネットや無線LANなどの送受信を行っている

> イーサーネット上を流れる電気信号（０と１が無限に流れているという意味でビットストリームと言います）を常に受信し続け、
> イーサーネットフレームの始まりと終わりを識別してイーサーネットフレームを取り出し、宛先MACアドレスを調べ、自ノード宛てかブロードキャスト宛ての場合、 MACヘッダー以外の部分を上位層（HTTPの例ではIP層）へ渡します。
- [ネットワーク - nicが行うイーサネットフレームの処理って、これですか？｜teratail](https://teratail.com/questions/75996)

cf.
- [MACアドレスとは(TCP/IP基礎)](http://ezxnet.com/network/entry4502/)

#### ネットワーク・インターフェースでよくあるやつ
- `en*` : Ethernet
- `ens*` ：有線の接続ポート。大体、最初は勝手にensの後に数字が振られる。
- `eth*` ：有線の接続ポート。上と内容は同じだが命名規則が変わった。こちらが旧。
- `lo*` ：ループバックのこと。実際はインターフェースとして存在しないが、テスト用などに仮想として必ずある。
- `virbr0` ：VMware使ってると出るらしい。
- `virbr0-nic` ：VMware使ってると出るらしい。nicはNetwork Interface Cardのこと

cf.
- [ネットワークインターフェイスの名前 - noyのブログ](http://noy.hatenablog.jp/entry/2017/02/27/163604)

## ちなみに
- ifconfig
- netstat
- arp
- route
は centos7 以降では非推奨となっていて、代わりに以下を使うことが推奨されているので注意しよう

- ifconfig	-> ip addr、ip -s link
- route	    -> ip route
- arp	    -> ip neigh
- netstat   -> ss

cf.
- [CentOS 7 以降では ifconfig、route、arp、netstat が非推奨 - eTuts+ Server Tutorial](https://server.etutsplus.com/centos-7-net-tools-vs-iproute2/)

## tcpdump

<details>

ネットワーク上を流れるデータを確認できるキャプチャツール
- 基本的にはTCP層以上のデータをキャプチャするが、オプション次第で、IP層のフレームやリンク層(MAC)のフレームもキャプチャできる

- ネットワーク通信の生データを収集し結果を出力(=パケットキャプチャ)してくれるCUIの解析ツール
- いつ、どこからどこへ、どんなフラグ(SYN,ACK,FIN等)のパケットが送られたか等が分かる。
- wireshark と同じことが(ほぼ)できる

### usage

``` shell
$ tcpdump port 80 and host 192.168.0.100 # 特定IPアドレスの80番ポートに関するトラフィックを見る(localhost (loopbackインターフェース)以外)

$ tcpdump -i lo0 port 8080 -nn # NIC を指定する場合(loopback NC の port 8080 のトラフィックを見る)
```

- -nn: ホスト名やポート番号をそのまま表示する
- -i: `tcpdump -i lo0 port 80` インタフェースを指定する。２枚刺しやブリッジとかしている場合は必須
- -e: 各ダンプ行ごとに、リンクレベルのヘッダを出力する
- -X: `tcpdump -X port 80` パケットの内容を 16進とASCIIで表示する。あまり使わない。ipヘッダまでを含めたパケットを見ることができる
- -XX: `tcpdump -XX port 80` イーサネット(リンクレベル)ヘッダまでを含めたフレームを16 進数と ASCIIで見ることができる
- -A: `tcpdump -A port 80` パケットの内容を ASCII で表示する。HTTP とか SMTP とかのキャプチャに便利
- -w file: `tcpdump -w {ファイル名} port 80` キャプチャ結果をファイルに出力する。出力されたファイルは Wireshark で開ける

mac の tcpdump は -i を指定しない場合 pkcap という疑似NC がデフォルトで選択されるようになっていて、すべての NC(loopback や ブリッジ以外)のトラフィックを見ることができる。
> On Darwin systems version 13 or later, when the interface is unspecified,
> tcpdump will use a pseudo interface to capture packets on a set of inter-
> faces determined by the kernel (excludes by default loopback and tunnel
> interfaces).
by man tcpdump

フォーマット

``` shell
time srcIP > dstIP: Flags [tcpflags], seq data-seqno, (or ack ackno), win window, (urg urgent), options [opts]
```

例
``` shell
通信時間           送信元IP              > 送信先IP                  フラグ     シーケンス番号  ウィンドウサイズ
16:58:46.899176 IP XXX.XXX.XXX.XXX.61372 > XXX.XXX.XXX.XXX.http-alt: Flags [S], seq 2816271339, win 65535, options [mss 1460,nop,wscale 4,nop,nop,TS val 33609312 ecr 0,sackOK,eol], length 0
```

- フラグ： TCPヘッダのフィールド内のコントロールフラグのこと cf. https://www.infraexpert.com/study/tcpip8.html

フラグの意味

- S (SYN)
- F (FIN)
- P (PUSH)
- R (RST)
- U (URG)
- W (ECN CWR)
- E (ECN-Echo)
- . (ACK)
- none (何もフラグがない場合)

> Tcpflags are some combination of S (SYN), F (FIN), P (PUSH), R (RST), U (URG), W (ECN CWR), E (ECN-Echo)  or `.'  (ACK), or `none' if no flags are set.

cf.
- [超絶初心者むけtcpdumpの使い方 - Qiita](https://qiita.com/tossh/items/4cd33693965ef231bd2a)
- [tcpdump の便利なオプション - Qiita](https://qiita.com/ngyuki/items/969d1efaddb68acb5313)
- [tcpdumpの使い方 - Qiita](https://qiita.com/aosho235/items/d87e0d69e89513d02a3f)
- mac 用の tcpdump についてや pktap(mac tcpdump で使用されている疑似ネットワークインターフェース) について調べてる [DSAS開発者の部屋:Mac OSX で vmnet が BIOCSETIF できなくてハマった話し](http://dsas.blog.klab.org/archives/52132993.html)
- [Man page of TCPDUMP](https://linuxjm.osdn.jp/html/tcpdump/man1/tcpdump.1.html)
- [tcpdump](http://support.tenasys.com/intimehelp_6_jp/util_tcpdump.html)


</details>

## tshark(wireshark)
NIC を指定して、そこを流れるデータを確認できるパケット(IP層のデータ)キャプチャツール
オプション次第で、MAC Addressフレームもキャプチャできる
tshark は wireshark の CUI版
mac だと `brew install wireshark && brew link wireshark` で tshark のみ入る

### usage
リアルタイムにパケットを表示する方法

``` shell
tshark -i lo0 -f "port 8080" # ネットワークインターフェース lo0(ループバックインターフェース) の port 8080 でフィルターしてキャプチャをする

tshark -i lo0 -Y "tcp.port==8080"
```

`tshark -i <インタフェース> -Y <絞り込みの条件> -n`

- -i: (ネットワーク)インタフェースを指定する
- -Y: 絞り込み条件を所定のフォーマットで指定する ex. `tcp.port==8080`
- -f: 絞り込みをする。こちらは、フォーマットで指定をせず grep のように使う
- -n: 名前解決をせずに数字のまま出力する
- -O ip,tcp: キャプチャを行う対象のプロトコルをコンマ区切りで指定する。 `tshark -G protocols` によって指定できるプロトコルを調べる事ができる
- -V: 要約でなく詳細を出力する

cf.
- [tsharkコマンドの使い方 - Qiita](https://qiita.com/hana_shin/items/0d997d9d9dd435727edf)
- [Wiresharkでパケットキャプチャしてみた - yagisukeのWebなブログ](http://yagisuke.hatenadiary.com/entry/2017/03/11/213513)
- [tsharkのインストールとフィルタ・自動停止オプションの使い方まとめ | OXY NOTES](https://oxynotes.com/?p=7969)
- [Wiresharkを使った通信監視（後編）――コマンドラインベースでのパケットキャプチャ | さくらのナレッジ](https://knowledge.sakura.ad.jp/6311/)
- [tshark - The Wireshark Network Analyzer 3.2.4](https://www.wireshark.org/docs/man-pages/tshark.html)

## 参考
- [ネットワークの基礎 : 51PM](http://51pm.blog.jp/archives/14371688.html)
- [2015年Webサーバアーキテクチャ序論 - ゆううきブログ](https://blog.yuuk.io/entry/2015-webserver-architecture)
- [Introduction | Learn You Some Erlang for Great Good!](https://www.ymotongpoo.com/works/lyse-ja/ja/26_buckets_of_sockets.html)
- [Working With Unix Processes — Learn the Fundamentals of Unix Programming in Ruby](https://www.jstorimer.com/products/working-with-unix-processes)
- [playground/books/working_with_tcp_sockets at master · ganmacs/playground](https://github.com/ganmacs/playground/tree/master/books/working_with_tcp_sockets)
- [Working with TCP Sockets 読書メモ 第6章 はじめてのクライアント／サーバ – Strings of Life](https://ryo511.info/archives/3809)
- [TCP/IP ソケットプログラミングの基礎を集中学習! Working with TCP sockets を読んでる | Futurismo](https://futurismo.biz/archives/2572)
- [「Working with TCP Sockets」を読んだ - Fire Engine](https://blog.tsurubee.tech/entry/2018/07/25/152514)qq
