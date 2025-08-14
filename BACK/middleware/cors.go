package middleware

import (
	"log"
	"net/http"
	"strings"
)

// CORS CORS 미들웨어
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Printf("CORS - 요청 Origin: %s", origin)

		// 192.168.0.x:3000 대역 모두 허용
		isAllowed := false
		if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
			isAllowed = true
		} else if strings.HasPrefix(origin, "http://192.168.0.") && strings.HasSuffix(origin, ":3000") {
			// 192.168.0.1:3000 ~ 192.168.0.254:3000 모두 허용
			isAllowed = true
		}

		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Printf("CORS - Origin 허용: %s", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			log.Printf("CORS - Origin 거부, 기본값 사용: %s", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// EnableCORS CORS 헤더 설정 (기존 코드와의 호환성을 위해)
func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
