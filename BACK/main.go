package main

import (
	"BACK/config"
	"BACK/handlers"
	"BACK/middleware"
	"BACK/service"
	"log"
	"net/http"
)

func main() {
	// 데이터베이스 초기화
	if err := config.InitDB(); err != nil {
		log.Fatal("데이터베이스 초기화 실패:", err)
	}

	// 서비스 생성
	authService := handlers.NewAuthService(config.DB)
	registerService := handlers.NewRegisterService(config.DB)
	jwtService := handlers.NewJWTService(config.GetDB())
	attendanceService := service.NewAttendanceService(config.DB)
	leaveService := handlers.NewLeaveService(config.DB)
	salaryService := handlers.NewSalaryService("./files/salary")
	scheduleService := handlers.NewScheduleService(config.DB)

	// 핸들러 생성
	authHandler := handlers.NewAuthHandler(authService)
	registerHandler := handlers.NewRegisterHandler(registerService)
	jwtHandler := handlers.NewJWTHandler(jwtService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	leaveHandler := handlers.NewLeaveHandler(leaveService)
	salaryHandler := handlers.NewSalaryHandler(salaryService)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)

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
	http.HandleFunc("/api/schedule", scheduleHandler.CreateSchedule)
	http.HandleFunc("/api/schedule/user", scheduleHandler.GetSchedulesByUser)
	http.HandleFunc("/api/schedule/update", scheduleHandler.UpdateSchedule)
	http.HandleFunc("/api/schedule/delete", scheduleHandler.DeleteSchedule)

	// 사용자 프로필 조회 라우트 추가
	http.HandleFunc("/api/user/profile", authService.GetUserProfile)

	// CORS 미들웨어 적용
	handler := middleware.CORS(http.DefaultServeMux)

	log.Println("서버가 포트 8081에서 시작되었습니다.")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
