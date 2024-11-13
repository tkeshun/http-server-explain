## 基本的なHTTPサーバー

### DefaultServeMux
http.HandlerFuncで特定のエンドポイントにリクエストが来たときの処理を設定する
http.ListenAndServeでサーバー起動
http.HandlerFuncを使うとデフォルトマルチプレクサにパスとリクエストを受けた際の処理が登録される

### ServeMux

HTTPを振り分けるためのマルチプレクサを自分で定義できる
構造体を１から定義せずとも、http.NewServeMuxでhttp.ServeMux構造体を取得できる
http.ServeMuxはhttp.Handlerインターフェースを満たす実装がされている

## カスタムハンドラー

http.Handlerインターフェースを満たすように実装することで、独自にServeHTTPメソッドを定義することができる
共通処理のカスタムハンドラーなどを作っておくことで、後述のミドルウェアチェーンなどで、ミドルウェアの再利用が可能

```
type AuthHandler struct {
	next http.Handler
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 簡単な認証チェック（例: ヘッダーにトークンがあるか）
	if r.Header.Get("Authorization") == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	h.next.ServeHTTP(w, r) // 認証後に次のハンドラーに制御を渡す
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are authorized!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/protected", protectedHandler)

	// カスタム認証ハンドラーを作成してServeMuxをラップ
	authHandler := &AuthHandler{next: mux}

	http.ListenAndServe(":8080", authHandler)
}

```


## ミドルウェアパターン

http.Handlerをラップすることで、ラップされるハンドラーの前後に処理を追加することができる
入れ子に組み合わせることで複数のミドルウェアを組み合わせられる

```
mux.Handle("/nested", loggingMiddleware(authMiddleware(http.HandlerFunc(nestedHandler))))
```

loggingMiddleware→authMiddleware→nestedHandlerの順に実行される

### ミドルウェアチェーン

複数のミドルウェアを組み合わせてリクエストが各ミドルウェアを順に通過するようにできる

```
func chainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}

// ミドルウェアの使用例
mux := http.NewServeMux()
mux.Handle("/", chainMiddleware(http.HandlerFunc(rootHandler), loggingMiddleware, authMiddleware))
```


## HTTPリクエストの処理


### リクエストヘッダーの取得

http.Request.Headerフィールドを使用するとリクエストヘッダーが取得できる

　
### コンテキストを使ったリクエストのキャンセル

*http.RequestのContextを使用してリクエストのキャンセルを検知できる。

```
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
```

## パスパラメータとワイルドカード

- 関連操作一覧

```
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
```

## グレースフルシャットダウン

シャットダウン時に必要な処理をしてからシャットダウンする
※コンテナで送られるシグナルについてdockerとかk8sのドキュメントにないか調べる
コンテナ環境では、プロセスがSIGTERMやSIGINTを受信したときに適切にシャットダウンできるように対応する必要がある

```
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
<-sigChan

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := server.Shutdown(ctx); err != nil {
	log.Fatalf("Server forced to shutdown: %v", err)
}
```

