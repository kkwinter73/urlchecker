package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheck_Success(t *testing.T) {
	// テスト用HTTPサーバーを立てる
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	result := Check(server.URL, 5*time.Second)

	if result.Err != nil {
		t.Errorf("エラーが発生: %v", result.Err)
	}
	if result.StatusCode != 200 {
		t.Errorf("期待値 200, 実際 %d", result.StatusCode)
	}
}

func TestCheck_InvalidURL(t *testing.T) {
	result := Check("http://this-domain-does-not-exist-xxxxx.invalid", 5*time.Second)

	if result.Err == nil {
		t.Error("エラーが返るべき")
	}
}

func TestCheckAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	urls := []string{server.URL, server.URL, server.URL}
	results := CheckAll(urls, 5*time.Second, 3)

	if len(results) != 3 {
		t.Errorf("結果が3件返るべき, 実際 %d件", len(results))
	}
	for _, r := range results {
		if r.StatusCode != 200 {
			t.Errorf("全て200を期待, 実際 %d", r.StatusCode)
		}
	}

}

func TestCheck_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := Check(server.URL, 1*time.Second)

	if result.Err == nil {
		t.Error("タイムアウトエラーが返るべき")
	}
}
