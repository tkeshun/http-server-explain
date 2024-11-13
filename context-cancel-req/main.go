package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	timer := time.After(5 * time.Second) // 5病後に実行
	for {
		select {
		case <-timer:
			fmt.Fprintln(w, "Request processed")
			return
		case <-ctx.Done(): // 5病経過前にリクエストをキャンセルされたら実行
			http.Error(w, "Request cancelled", http.StatusRequestTimeout)
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	// マルチプレクサ(handler)を登録する必要がある
	// 登録する場合、第一引数のポート番号と第２引数のhandlerがServer構造体に登録された上で、ServerのListenAndServeメソッドが実行される
	http.ListenAndServe(":8080", mux)
}
