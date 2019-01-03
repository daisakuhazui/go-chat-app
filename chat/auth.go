package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証
		log.Printf("You have to log in. Redirecting login page.")
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 何らかのエラーが発生
		panic(err.Error())
	} else {
		// 成功 ラップされたハンドラを呼び出します
	}
	h.next.ServeHTTP(w, r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler はサードパーティーへのログイン処理を待ち受けます
// パスの形式: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		// login アクション
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("Fail to get authentication provider:", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("ERROR occured in calling GetBeginAuthURL:", provider, "-", err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		// callback アクション
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("Fail to get authentication provider:", provider, "-", err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("Incomplete authentication", provider, "-", err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("Fail to get user", provider, "-", err)
		}
		autheCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: autheCookieValue,
			Path:  "/",
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		// サポートされていないアクション
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Unsupported Action: %s", action)
	}
}
