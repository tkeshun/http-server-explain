package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", handler) // DefaultServeMuxにパスと対応する処理が登録される
	// ListenAndServeでマルチプレクサが登録されない場合、DefaultServeMuxが使用される
	// ライブラリ内のコード。/usr/local/go/src/net/http/server.goにhandler判別ロジックがある
	// ↓該当コード
	// func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	// 	handler := sh.srv.Handler
	// 	if handler == nil {
	// 		handler = DefaultServeMux
	// 	}
	// 	if !sh.srv.DisableGeneralOptionsHandler && req.RequestURI == "*" && req.Method == "OPTIONS" {
	// 		handler = globalOptionsHandler{}
	// 	}

	// 	handler.ServeHTTP(rw, req)
	// }
	http.ListenAndServe(":8080", nil)
}
