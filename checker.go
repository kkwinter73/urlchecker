package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

// String は結果を人間が読める形に整形する
func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("[NG] %s -> error: %v (%.2fs)", r.URL, r.Err, r.Duration.Seconds())
	}
	return fmt.Sprintf("[OK] %s -> %d (%.2fs)", r.URL, r.StatusCode, r.Duration.Seconds())
}

func Check(url string, timeout time.Duration) Result {
	start := time.Now()

	// タイムアウト付きのコンテキストを作る
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// contextを使ってリクエストを組み立てる
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{URL: url, Err: err}
	}

	// リクエストを送信する
	resp, err := http.DefaultClient.Do(req)

	duration := time.Since(start)

	if err != nil {
		return Result{URL: url, Duration: duration, Err: err}
	}

	defer resp.Body.Close()

	return Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}

func CheckAll(urls []string, timeout time.Duration) []Result {
	results := make([]Result, len(urls)) // makeで長さ固定のスライスを作る。最初から長さを固定しているのはあとからアクセスしやすいようにするため

	// channelの作成 - 平衡実行されるgoroutine間を接続するパイプ的な
	// 並行実行している関数から値を受信する
	// （あるgoroutineから別のgoroutineに値を渡す）
	// make(chan 型)で新しいチャネルを作成できる
	// channel <- 構文で、チャネルへ値を 送信 します。
	// <-channel 構文で、チャネルから値を 受信 します
	ch := make(chan struct { // 匿名構造体：GOは型システムが厳格なため構造体で書く必要がある
		index  int
		result Result
	}, len(urls)) // バッファ容量 len(urls）のチャネル：今回はN個おくることがわかってるので、バッファをその分もたせることで、全goroutineが受信をまたずにサクッと終了できる

	for i, url := range urls {
		// goroutine の起動
		go func(index int, u string) {
			ch <- struct {
				index  int
				result Result
			}{index, Check(u, timeout)} // Check(U)でurlにアクセスしてResultを得る
		}(i, url)
	}

	// goroutineは並行実行なので完了順がバラバラ
	// URLの元の順番に結果を並べたいならindexを一緒に送って、添え字で正しい場所に格納する
	for range urls {
		r := <-ch
		results[r.index] = r.result
	}

	return results
}
