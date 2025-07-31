package main

import (
	"BACK/config"
	"BACK/handlers"
	"BACK/middleware"
	"BACK/utils"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func main() {
	// JWT 시크릿 초기화
	utils.InitJWTSecret()
	config.InitDB()
	defer config.DB.Close()

	publicMux := http.NewServeMux()
	privateMux := http.NewServeMux()

	// 공개 API (인증 없이 접근 가능)
	publicMux.HandleFunc("/api/register", handlers.HandleRegister)
	publicMux.HandleFunc("/api/login", handlers.HandleLogin)
	publicMux.HandleFunc("/api/jwt/login", handlers.HandleJWTLogin)
	publicMux.HandleFunc("/api/jwt/register", handlers.HandleRegister)
	publicMux.HandleFunc("/api/refresh", handlers.HandleRefreshToken)

	// 인증이 필요한 API
	privateMux.HandleFunc("/api/logout", handlers.HandleLogout)
	privateMux.HandleFunc("/api/user/profile", handlers.GetCurrentUser)
	privateMux.HandleFunc("/api/jwt/logout", handlers.HandleJWTLogout)
	privateMux.HandleFunc("/api/jwt/refresh", handlers.HandleTokenRefresh)
	privateMux.HandleFunc("/api/jwt/verify", handlers.HandleVerifyToken)
	privateMux.HandleFunc("/api/jwt/user/profile", func(w http.ResponseWriter, r *http.Request) {
		userID, username, role, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
		response := map[string]interface{}{
			"message": "JWT 기반 인증된 사용자의 프로필 정보입니다",
			"user": map[string]interface{}{
				"user_id":  userID,
				"username": username,
				"role":     role,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	privateMux.HandleFunc("/api/jwt/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		userID, username, role, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
		response := map[string]interface{}{
			"message": "관리자 대시보드에 접근했습니다",
			"user": map[string]interface{}{
				"user_id":  userID,
				"username": username,
				"role":     role,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	privateMux.HandleFunc("/api/schedule", handlers.AddSchedule)
	privateMux.HandleFunc("/api/schedule/", handlers.HandleScheduleOperations)
	privateMux.HandleFunc("/api/schedules", handlers.GetSchedulesByUser)

	// 출퇴근 관련 API
	privateMux.HandleFunc("/api/attendance/checkin", handlers.CheckIn)
	privateMux.HandleFunc("/api/attendance/checkout", handlers.CheckOut)
	privateMux.HandleFunc("/api/attendance/records/", handlers.GetAttendanceRecords)
	privateMux.HandleFunc("/api/attendance/today/", handlers.GetTodayAttendance)
	privateMux.HandleFunc("/api/attendance/date/", handlers.GetAttendanceByDate)

	// 연차 관련 API
	privateMux.HandleFunc("/api/leave/status", handlers.GetLeaveStatus)
	privateMux.HandleFunc("/api/leave/policy", handlers.GetLeavePolicy)
	privateMux.HandleFunc("/api/leave/calculate", handlers.CalculateLeaveForDate)
	privateMux.HandleFunc("/api/leave/update-all", handlers.UpdateAllUsersLeaveStatus)

	// 급여명세서 관련 API
	privateMux.HandleFunc("/api/salary/statement", handlers.GetSalaryStatement)
	privateMux.HandleFunc("/api/salary/list", handlers.GetSalaryStatementList)
	privateMux.HandleFunc("/api/dual-auth/profile", func(w http.ResponseWriter, r *http.Request) {
		userID, username, role, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
		response := map[string]interface{}{
			"message": "세션+JWT 모두 인증된 사용자의 프로필 정보입니다",
			"user": map[string]interface{}{
				"user_id":  userID,
				"username": username,
				"role":     role,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("요청 경로: %s, 메서드: %s", r.URL.Path, r.Method)

		// 공개 API (인증 없이 접근 가능)
		if r.URL.Path == "/api/login" || r.URL.Path == "/api/register" ||
			r.URL.Path == "/api/jwt/login" || r.URL.Path == "/api/jwt/register" ||
			r.URL.Path == "/api/refresh" {
			log.Printf("공개 API 접근: %s", r.URL.Path)
			middleware.EnableCORS(publicMux).ServeHTTP(w, r)
			return
		}

		// 인증이 필요한 API
		if strings.HasPrefix(r.URL.Path, "/api/") {
			log.Printf("인증 필요 API 접근: %s", r.URL.Path)

			// 급여명세서 API 특별 로그
			if strings.Contains(r.URL.Path, "/api/salary/") {
				log.Printf("[SALARY API] 급여명세서 API 요청: %s", r.URL.Path)
				log.Printf("[SALARY API] 요청 파라미터: %s", r.URL.RawQuery)
			}

			middleware.EnableCORS(middleware.DualAuthMiddleware(privateMux)).ServeHTTP(w, r)
			return
		}

		// 정적 파일
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
	})

	log.Println("API 서버가 http://localhost:8080 에서 실행 중입니다...")
	log.Println("/api/login, /api/register, /api/jwt/login, /api/jwt/register, /api/refresh는 인증 없이 접근 가능합니다.")
	log.Println("그 외 모든 /api/* 엔드포인트는 세션+JWT 복합 인증이 필요합니다.")
	log.Println("로그인(/api/login) 시 세션+JWT가 모두 발급됩니다.")
	log.Println("프론트엔드에서 세션 쿠키와 JWT 토큰을 모두 전송해야 API 접근이 가능합니다.")

	log.Fatal(http.ListenAndServe(":8080", finalHandler))
}
