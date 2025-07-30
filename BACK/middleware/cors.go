package middleware

import (
	"net/http"
)

// EnableCORS는 CORS 헤더를 설정하는 미들웨어 함수입니다
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS 헤더 설정
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Vue 개발 서버 포트
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // 중요: 쿠키 허용

		// OPTIONS 요청 처리
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//다른 도메인(출처)에서 요청을 보내도 서버가 허용할지 말지 결정하는 정책.
