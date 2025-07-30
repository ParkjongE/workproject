package config

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("secret-key-12345"))

func init() {
	// 세션 스토어 설정
	Store.Options = &sessions.Options{
		Path:     "/",                  // 모든 경로에서 접근 가능
		MaxAge:   3600 * 24,            // 24시간 유지
		HttpOnly: true,                 // JavaScript에서 접근 불가
		Secure:   false,                // HTTPS가 아닌 경우 false
		SameSite: http.SameSiteLaxMode, // CSRF 방지
	}
}

// GetSessionStore 함수를 추가하여 다른 패키지에서 세션 스토어에 접근할 수 있도록 함
func GetSessionStore() *sessions.CookieStore {
	return Store
}
