package middleware

import (
	"log"
	"net/http"

	"BACK/config"

	"github.com/gorilla/sessions"
)

// SessionAuth는 세션 기반 인증을 처리하는 미들웨어 함수입니다
func SessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 인증이 필요 없는 경로는 바로 통과
		if r.URL.Path == "/api/login" || r.URL.Path == "/api/register" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// 세션 가져오기
		session, err := config.GetSessionStore().Get(r, "user-session")
		if err != nil {
			log.Printf("세션 가져오기 실패: %v", err)
			http.Error(w, "세션 오류", http.StatusInternalServerError)
			return
		}

		// 세션에서 사용자 정보 확인
		authenticated, ok := session.Values["authenticated"].(bool)
		if !ok || !authenticated {
			log.Println("세션 인증되지 않음")
			http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
			return
		}

		userID, ok := session.Values["user_id"].(int)
		if !ok {
			log.Println("세션에 사용자 ID가 없습니다")
			http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
			return
		}

		username, ok := session.Values["username"].(string)
		if !ok {
			log.Println("세션에 사용자명이 없습니다")
			http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
			return
		}

		role, ok := session.Values["role"].(string)
		if !ok {
			log.Println("세션에 역할이 없습니다")
			http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
			return
		}

		log.Printf("인증된 사용자: %s (ID: %d, 역할: %s)", username, userID, role)

		// 다음 핸들러 호출
		next.ServeHTTP(w, r)
	})
}

// 세션 스토어를 외부에서 접근할 수 있도록 Getter 함수 제공
func GetSessionStore() *sessions.CookieStore {
	return config.GetSessionStore()
}
