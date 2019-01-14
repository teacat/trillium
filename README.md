# Trillium [![GoDoc](https://godoc.org/github.com/teacat/trillium?status.svg)](https://godoc.org/github.com/teacat/trillium) [![Coverage Status](https://coveralls.io/repos/github/teacat/trillium/badge.svg?branch=master)](https://coveralls.io/github/teacat/trillium?branch=master) [![Build Status](https://travis-ci.org/teacat/trillium.svg?branch=master)](https://travis-ci.org/teacat/trillium) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/trillium)](https://goreportcard.com/report/github.com/teacat/trillium)

Trillium 是一個基於 TeaCat 所需而提出的分散式不重複唯一編號演算規則，其方式與 Twitter 所設計的 [Snowflake 雪花編號](https://developer.twitter.com/en/docs/basics/twitter-ids.html)類似。

* 不須要中心化服務，任何一個服務都能獨立產生唯一編號。
* 基於時間順序而排定的編號，能夠更方便排序。
* 無流水號問題而能被得知總體數量或是資料被電腦程式爬取。
* 單個服務每秒可以產生 100,000 個唯一編號；若額度耗盡，將會暫停動作並延遲到下一秒才繼續配發唯一編號。
* 最高可高達 99,999 個服務同時使用，或是 9,999 個（隨機碰撞較安全的範圍）。
* 可用編號時間配發時間高達 292 年。

簡單來說：Trillium 最可以同時執行在至少 9,999 個服務中；每個服務每秒最高可以處理 100,000 個請求（略估為每毫秒 100 個請求）；照這個方式下去，編號將能持續提供到 292 年後。

## 構造

Trillium 只能執行在 64 位元的電腦中，因為其編號長度高達 20 字元寬度。

```txt
   已過時間      機器隨機編號    流水編號
+------------+-------------+---------+
| 1547491194 |    61835    |  01824  |  = "15474911946183501824"
+------------+-------------+---------+
    10 字元        5 字元      5 字元
```

## 問與答

**問：如何確保不同服務不會產生重複的編號？**
答：當 Trillium 被建立時，會自動產生長度為 5 的隨機數字來作為服務的唯一辨識號碼並避免與其他正在執行的服務產生出相同的編號。

**問：看起來編號中有隨機要素存在，那麼為什麼這個編號還可以按照時間排序？**
答：因為 Trillium 以時間（秒數）為主並將其擺放在開頭，所以就算後面的要素導致亂序，但因為最前面是以時間開頭為基準，所以這個唯一編號仍能保持有序。但也因為 Trillium 是以秒數為基準單位，因此排序僅能精準到「秒」，無法算出到「毫秒」的排序（_不過也夠大部分場合使用了_）。

**問：用數字編號作為唯一編號，不是會很容易被流水號攻擊、讓別人知道資料總筆數量，或是被機器人爬取資料嗎？**
答：不會。Trillium 每秒就會更改一次唯一編號的規則（開頭秒數的異動），而且加上不同服務有著不同的隨機編號，因此無法透過遞增的方式來爬取、或是得知資料的總筆數量為何。

## 可參考文件

[分布式系统中 Unique ID 的生成方法](https://darktea.github.io/notes/2013/12/08/Unique-ID)
[扯扯ID - 掘金](https://juejin.im/post/593d0821128fe1006ae47e3c)