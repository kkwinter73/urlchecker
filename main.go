package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	urlsFlag := flag.String("urls", "", "チェックするURL（カンマ区切り）")
	timeoutFlag := flag.Int("timeout", 10, "タイムアウト秒数")
	concurrencyFlag := flag.Int("concurrency", 3, "同時実行数")

	flag.Parse()

	if *urlsFlag == "" {
		fmt.Println("使い方: url-checker --urls=https://example.com,https://google.com")
		os.Exit(1)
	}

	if *concurrencyFlag <= 0 {
		fmt.Println("１以上の入力が必要")
		os.Exit(1)
	}

	urls := strings.Split(*urlsFlag, ",")

	start := time.Now()

	timeout := time.Duration(*timeoutFlag) * time.Second

	results := CheckAll(urls, timeout, *concurrencyFlag)
	total := time.Since(start)

	for _, r := range results {
		fmt.Println(r)
	}
	fmt.Printf("\n合計時間: %.2fs\n", total.Seconds())
}
