青空文庫の作品などをノベルゲーム風に読めるようにする。

![images](./top.png)

- フォントは http://jikasei.me/font/jf-dotfont/ から拝借した
- 本文は夏目漱石『坊っちゃん』 https://www.aozora.gr.jp/cards/000148/files/752_14964.html から拝借した

## 青空文庫のルビを削除する方法

青空文庫のテキストにはルビが含まれていて、そのままコピペすると文章がおかしくなる。

https://qiita.com/kanaxx/items/6d6b0d680185d6af9b05 で紹介されている方法。

ブラウザのコンソールで↓を実行する。

```js
$('rt').hide();
```

ページによってはなぜかできないこともあるので、↓で。
たとえばhttps://www.aozora.gr.jp/cards/000148/files/1747_14970.html

```js
document.querySelectorAll("rt").forEach(el => el.remove());
```
