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
	flag.Parse()

	if *urlsFlag == "" {
		fmt.Println("使い方: url-checker --urls=https://example.com,https://google.com")
		os.Exit(1)
	}

	urls := strings.Split(*urlsFlag, ",")

	start := time.Now()
	results := CheckAll(urls)
	total := time.Since(start)

	for _, r := range results {
		fmt.Println(r)
	}
	fmt.Printf("\n合計時間: %.2fs\n", total.Seconds())
}
