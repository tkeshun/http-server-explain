package main

import "net/http"

type CustomHandler struct{} // anyだとエラーになる

// http.Handlerのインターフェースを実装する
func (h *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Custom handler\n"))
}

func main() {
	handler := &CustomHandler{}
	// DefaultServeMuxにカスタムハンドラーを登録する
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
