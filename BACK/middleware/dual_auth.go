package middleware

import (
	"context"
	"log"
	"net/http"

	"BACK/config"
	"BACK/utills"
)

// DualAuthMiddleware는 세션과 JWT가 모두 유효해야 인증 성공
func DualAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[DualAuth] 인증 시작: %s", r.URL.Path)

		// 1. 세션 인증
		session, err := config.GetSessionStore().Get(r, "user-session")
		if err != nil {
			log.Printf("[DualAuth] 세션 가져오기 실패: %v", err)
			http.Error(w, "세션 오류", http.StatusUnauthorized)
			return
		}

		authenticated, ok := session.Values["authenticated"].(bool)
		userID, okID := session.Values["user_id"].(int)
		username, okName := session.Values["username"].(string)
		role, okRole := session.Values["role"].(string)

		log.Printf("[DualAuth] 세션 값: authenticated=%v, userID=%d, username=%s, role=%s",
			authenticated, userID, username, role)

		if !ok || !authenticated || !okID || !okName || !okRole {
			log.Printf("[DualAuth] 세션 인증 실패: ok=%v, authenticated=%v, okID=%v, okName=%v, okRole=%v",
				ok, authenticated, okID, okName, okRole)
			http.Error(w, "세션 인증 필요", http.StatusUnauthorized)
			return
		}

		log.Printf("[DualAuth] 세션 인증 성공")

		// 2. JWT 인증
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("[DualAuth] Authorization 헤더 없음")
			http.Error(w, "JWT 인증 필요", http.StatusUnauthorized)
			return
		}

		tokenString, err := utills.ExtractTokenFromHeader(authHeader)
		if err != nil {
			log.Printf("[DualAuth] JWT 추출 실패: %v", err)
			http.Error(w, "JWT 인증 필요", http.StatusUnauthorized)
			return
		}

		claims, err := utills.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("[DualAuth] JWT 검증 실패: %v", err)
			http.Error(w, "JWT 인증 필요", http.StatusUnauthorized)
			return
		}

		log.Printf("[DualAuth] JWT 인증 성공: userID=%d, username=%s, role=%s",
			claims.UserID, claims.Userid, claims.Role)

		// 3. 세션 정보와 JWT 정보 일치 여부(선택)
		log.Printf("[DualAuth] 세션과 JWT 비교: 세션(userID=%d, username=%s, role=%s) vs JWT(userID=%d, username=%s, role=%s)",
			userID, username, role, claims.UserID, claims.Userid, claims.Role)

		if userID != claims.UserID || username != claims.Userid || role != claims.Role {
			log.Printf("[DualAuth] 세션과 JWT 사용자 정보 불일치: 세션(userID=%d, username=%s, role=%s) vs JWT(userID=%d, username=%s, role=%s)",
				userID, username, role, claims.UserID, claims.Userid, claims.Role)
			http.Error(w, "인증 정보 불일치", http.StatusUnauthorized)
			return
		}

		log.Printf("[DualAuth] 모든 인증 성공 - 요청 처리 진행")

		// 4. Context에 사용자 정보 저장
		ctx := context.WithValue(r.Context(), UserContextKeyID, userID)
		ctx = context.WithValue(ctx, UserContextKeyUsername, username)
		ctx = context.WithValue(ctx, UserContextKeyRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
