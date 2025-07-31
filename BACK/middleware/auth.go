package middleware

import (
	"BACK/utils"
	"context"
	"log"
	"net/http"
)

// UserContextKey는 컨텍스트에서 사용자 정보를 저장할 때 사용하는 키입니다
type UserContextKey string

const (
	UserContextKeyID       UserContextKey = "user_id"
	UserContextKeyUsername UserContextKey = "username"
	UserContextKeyRole     UserContextKey = "role"
)

// JWTAuthMiddleware는 JWT 토큰을 검증하는 미들웨어입니다
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 인증이 필요 없는 경로는 바로 통과
		if r.URL.Path == "/api/jwt/login" || r.URL.Path == "/api/jwt/register" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Authorization 헤더 확인
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Authorization 헤더가 없습니다")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Bearer 토큰 추출
		tokenString, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			log.Printf("토큰 추출 실패: %v", err)
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// JWT 토큰 검증
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("JWT 토큰 검증 실패: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 컨텍스트에 사용자 정보 저장
		ctx := context.WithValue(r.Context(), UserContextKeyID, claims.UserID)
		ctx = context.WithValue(ctx, UserContextKeyUsername, claims.Userid)
		ctx = context.WithValue(ctx, UserContextKeyRole, claims.Role)

		// 다음 핸들러 호출
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext는 컨텍스트에서 사용자 정보를 가져옵니다
func GetUserFromContext(ctx context.Context) (int, string, string, bool) {
	userID, ok1 := ctx.Value(UserContextKeyID).(int)
	username, ok2 := ctx.Value(UserContextKeyUsername).(string)
	role, ok3 := ctx.Value(UserContextKeyRole).(string)

	if !ok1 || !ok2 || !ok3 {
		return 0, "", "", false
	}

	return userID, username, role, true
}

// RequireRole은 특정 역할을 가진 사용자만 접근할 수 있도록 하는 미들웨어입니다
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 사용자 정보 가져오기
			userID, username, role, ok := GetUserFromContext(r.Context())
			if !ok {
				log.Println("컨텍스트에서 사용자 정보를 찾을 수 없습니다")
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}

			// 역할 확인
			if role != requiredRole {
				log.Printf("권한 부족: 사용자 %s(%d)가 %s 역할을 요구하지만 %s 역할을 가지고 있습니다",
					username, userID, requiredRole, role)
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// 다음 핸들러 호출
			next.ServeHTTP(w, r)
		})
	}
}

// OptionalJWTAuthMiddleware는 JWT 토큰이 있으면 검증하고, 없어도 통과시키는 미들웨어입니다
func OptionalJWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorization 헤더 확인
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// 토큰이 없어도 통과
			next.ServeHTTP(w, r)
			return
		}

		// Bearer 토큰 추출
		tokenString, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// 토큰 형식이 잘못되어도 통과 (선택적 인증이므로)
			log.Printf("토큰 형식 오류 (선택적 인증): %v", err)
			next.ServeHTTP(w, r)
			return
		}

		// JWT 토큰 검증
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			// 토큰이 유효하지 않아도 통과 (선택적 인증이므로)
			log.Printf("토큰 검증 실패 (선택적 인증): %v", err)
			next.ServeHTTP(w, r)
			return
		}

		// 컨텍스트에 사용자 정보 저장
		ctx := context.WithValue(r.Context(), UserContextKeyID, claims.UserID)
		ctx = context.WithValue(ctx, UserContextKeyUsername, claims.Userid)
		ctx = context.WithValue(ctx, UserContextKeyRole, claims.Role)

		// 다음 핸들러 호출

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORS와 JWT 인증을 함께 처리하는 미들웨어
func CORSWithJWTAuth(next http.Handler) http.Handler {
	return EnableCORS(JWTAuthMiddleware(next))
}

// CORS와 선택적 JWT 인증을 함께 처리하는 미들웨어
func CORSWithOptionalJWTAuth(next http.Handler) http.Handler {
	return EnableCORS(OptionalJWTAuthMiddleware(next))
}
