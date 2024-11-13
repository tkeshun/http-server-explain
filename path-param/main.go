package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// Method指定
	mux.HandleFunc("GET /method", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(r.Method)
		fmt.Fprintf(w, "Method: %s\n", r.Method)
	})

	// パスパラメータ
	mux.HandleFunc("/pathvalue/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Println(id)
		fmt.Fprintf(w, "ID: %s\n", id)
	})

	// ワイルドカード
	mux.HandleFunc("/wild/{wild...}", func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("wild")
		fmt.Println(path)
		fmt.Fprintf(w, "Requested path: %s\n", path)
	})

	// 末尾スラッシュの完全一致のみ許容する
	mux.HandleFunc("/perfect/{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("末尾スラッシュの完全一致のみ許容する")
		fmt.Fprintf(w, "完全一致\n")
	})

	// マルチプレクサ(handler)を登録する必要がある
	// 登録する場合、第一引数のポート番号と第２引数のhandlerがServer構造体に登録された上で、ServerのListenAndServeメソッドが実行される
	http.ListenAndServe(":8080", mux)
}
