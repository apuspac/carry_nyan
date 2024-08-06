oauth tokenを取ってくるやつ。
scopeをとりあえずばかすか入れてるけど、 ほんとは良くないと思います。
完全個人用で絶対流出させない自信があれば大丈夫。なのかなぁ。

scopeと、client idを設定。
client idは、twitch_developerでアプリ登録したときに取得できる。
scopeは、doc見よう。

```
bun run oauth_http.ts
```

認証ボタン押すと、twitchの許可しますか？に飛ぶので、許可すると、
urlにaccess tokenが乗っかってくるので、それをコピーする。
