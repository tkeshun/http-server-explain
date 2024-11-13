package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func jsonRequestHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストのBodyを閉じるためのdefer
	defer r.Body.Close()

	// Bodyを読み込んで構造体にデコード
	var data RequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// デコード結果を出力
	fmt.Fprintf(w, "Received JSON: Name=%s, Email=%s", data.Name, data.Email)
}

type ResponseData struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func jsonResponseHandler(w http.ResponseWriter, r *http.Request) {
	// レスポンスデータを作成
	response := ResponseData{
		Message: "Hello, JSON World!",
		Status:  200,
	}

	// ヘッダーのContent-Typeを設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// JSONデータをエンコードしてレスポンスBodyに書き込む
	json.NewEncoder(w).Encode(response)
}

func jsonHandlerWithError(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// リクエストサイズを制限
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MBの制限

	var data RequestData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // 未知のフィールドを拒否

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// 正常時のレスポンス
	fmt.Fprintf(w, "Received JSON: Name=%s, Email=%s", data.Name, data.Email)
}

func main() {
	http.HandleFunc("/request", jsonRequestHandler)
	http.HandleFunc("/response", jsonResponseHandler)
	http.ListenAndServe(":8080", nil)
}
