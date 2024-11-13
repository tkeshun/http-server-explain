package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	// マルチプレクサ(handler)を登録する必要がある
	// 登録する場合、第一引数のポート番号と第２引数のhandlerがServer構造体に登録された上で、ServerのListenAndServeメソッドが実行される
	http.ListenAndServe(":8080", mux)
}
