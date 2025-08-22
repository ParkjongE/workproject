package main

import (
	"BACK/config"
	"BACK/handler"
	"BACK/middleware"
	"BACK/repository"
	"BACK/service"
	"log"
	"net/http"
	"strings"
)

func main() {
	// 데이터베이스 초기화
	if err := config.InitDB(); err != nil {
		log.Fatal("데이터베이스 초기화 실패:", err)
	}

	// Repository 생성
	attendanceRepo := repository.NewAttendanceRepository(config.DB)
	scheduleRepo := repository.NewScheduleRepository(config.DB)
	registerRepo := repository.NewRegisterRepository(config.DB)
	leaveRepo := repository.NewLeaveRepository(config.DB)
	scheduleSharingRepo := repository.NewScheduleSharingRepository(config.GetDB()) // 수정: GetDB() 사용

	// 서비스 생성
	authService := handler.NewAuthService(config.DB)
	registerService := service.NewRegisterService(registerRepo)
	jwtService := handler.NewJWTService(config.GetDB())
	attendanceService := service.NewAttendanceService(attendanceRepo)
	leaveService := service.NewLeaveService(leaveRepo)
	salaryService := service.NewSalaryService("./files/salary")
	scheduleService := service.NewScheduleService(scheduleRepo)
	scheduleSharingService := service.NewScheduleSharingService(scheduleSharingRepo) // 추가

	// 핸들러 생성
	authHandler := handler.NewAuthHandler(authService)
	registerHandler := handler.NewRegisterHandler(registerService)
	jwtHandler := handler.NewJWTHandler(jwtService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	leaveHandler := handler.NewLeaveHandler(leaveService)
	salaryHandler := handler.NewSalaryHandler(salaryService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	scheduleSharingHandler := handler.NewScheduleSharingHandler(scheduleSharingService) // 추가

	// 라우터 설정
	http.HandleFunc("/api/register", registerHandler.HandleRegister)
	http.HandleFunc("/api/login", authHandler.HandleLogin)
	http.HandleFunc("/api/jwt/login", jwtHandler.HandleJWTLogin)
	http.HandleFunc("/api/jwt/logout", jwtHandler.HandleJWTLogout)
	http.HandleFunc("/api/jwt/refresh", jwtHandler.HandleTokenRefresh)
	http.HandleFunc("/api/jwt/verify", jwtHandler.HandleVerifyToken)

	// 출근 관련
	http.HandleFunc("/api/attendance/history", attendanceHandler.GetAttendanceHistory)  // GET: 출근 이력
	http.HandleFunc("/api/attendance/date/", attendanceHandler.GetAttendanceByDate)     // GET: 특정 날짜 (/ 포함)
	http.HandleFunc("/api/attendance/records/", attendanceHandler.GetAttendanceRecords) // GET: 사용자 이력
	http.HandleFunc("/api/attendance/checkin", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("=== ROUTER LOG === 요청 수신: %s %s", r.Method, r.URL.Path)
		attendanceHandler.CheckIn(w, r)
	}) // POST: 출근
	http.HandleFunc("/api/attendance/checkout", attendanceHandler.CheckOut)       // POST: 퇴근
	http.HandleFunc("/api/attendance/record", attendanceHandler.RecordAttendance) // POST: 출근 기록 (호환성)

	// 연차 관련
	http.HandleFunc("/api/leave/status", leaveHandler.GetUserLeaveStatus)
	http.HandleFunc("/api/leave/calculate", leaveHandler.CalculateLeaveForDate)
	http.HandleFunc("/api/leave/update-all", leaveHandler.UpdateAllUsersLeaveStatus)

	// 급여 관련
	http.HandleFunc("/api/salary/statement", salaryHandler.GetSalaryStatement)
	http.HandleFunc("/api/salary/list", salaryHandler.GetSalaryStatementList)

	// 스케줄 관련
	http.HandleFunc("/api/schedule", scheduleHandler.CreateSchedule)        // POST: 일정 생성
	http.HandleFunc("/api/schedules", scheduleHandler.GetSchedulesByUser)   // GET: 일정 조회
	http.HandleFunc("/api/schedule/update", scheduleHandler.UpdateSchedule) // PUT: 일정 수정
	http.HandleFunc("/api/schedule/delete", scheduleHandler.DeleteSchedule) // DELETE: 일정 삭제

	// 스케줄 공유 관련 라우트 추가
	http.HandleFunc("/api/schedule-sharing/users/all", scheduleSharingHandler.GetAllDepartmentUsers)
	http.HandleFunc("/api/schedule-sharing/users/department/", scheduleSharingHandler.GetUsersByDepartment)
	http.HandleFunc("/api/schedule-sharing/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			path := strings.TrimPrefix(r.URL.Path, "/api/schedule-sharing/users/")
			if path == "" {
				http.Error(w, "User ID required", http.StatusBadRequest)
				return
			}

			// /api/schedule-sharing/users/{userId}/schedules
			if strings.HasSuffix(path, "/schedules") {
				userIDStr := strings.TrimSuffix(path, "/schedules")
				if userIDStr == "" {
					http.Error(w, "User ID required", http.StatusBadRequest)
					return
				}
				scheduleSharingHandler.GetUserSchedulesByID(w, r)
			} else {
				// /api/schedule-sharing/users/{userId}
				scheduleSharingHandler.GetUserSchedulesByID(w, r) // GetUserWithSchedulesByID → GetUserSchedulesByID
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 사용자 프로필 조회 라우트 추가
	http.HandleFunc("/api/user/profile", authService.GetUserProfile)

	// 요청 로깅 미들웨어 추가
	loggingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🚀 요청: %s %s", r.Method, r.URL.String())
		http.DefaultServeMux.ServeHTTP(w, r)
	})

	// CORS 미들웨어 적용
	handler := middleware.CORS(loggingHandler)

	log.Println("서버가 포트 8081에서 시작되었습니다.")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
