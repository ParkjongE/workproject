package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"BACK/config"
	"BACK/utils"
)

// GetLeaveStatus 사용자 연차 현황 조회
func GetLeaveStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GetLeaveStatus] 요청 시작")

	// URL 파라미터에서 user_id 추출
	userIDStr := r.URL.Query().Get("user_id")
	log.Printf("[GetLeaveStatus] user_id 파라미터: %s", userIDStr)

	if userIDStr == "" {
		log.Printf("[GetLeaveStatus] user_id 파라미터 없음")
		http.Error(w, "user_id 파라미터가 필요합니다", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("[GetLeaveStatus] user_id 변환 오류: %v", err)
		http.Error(w, "잘못된 user_id입니다", http.StatusBadRequest)
		return
	}
	log.Printf("[GetLeaveStatus] 변환된 userID: %d", userID)

	// 현재 연도 사용
	year := time.Now().Year()
	log.Printf("[GetLeaveStatus] 현재 연도: %d", year)

	// 사용자 정보 조회 (입사일 포함)
	var hireDateStr string
	log.Printf("[GetLeaveStatus] 사용자 정보 조회 시작: userID=%d", userID)
	log.Printf("[GetLeaveStatus] DB 쿼리 실행: SELECT hire_date FROM users WHERE id = %d", userID)

	err = config.DB.QueryRow("SELECT hire_date FROM users WHERE id = ?", userID).Scan(&hireDateStr)
	if err != nil {
		log.Printf("[GetLeaveStatus] 사용자 정보 조회 오류: %v", err)
		log.Printf("[GetLeaveStatus] 오류 타입: %T", err)
		if err == sql.ErrNoRows {
			log.Printf("[GetLeaveStatus] 사용자를 찾을 수 없음: userID=%d", userID)
			http.Error(w, "사용자를 찾을 수 없습니다", http.StatusNotFound)
			return
		}
		log.Printf("[GetLeaveStatus] 데이터베이스 오류 발생")
		http.Error(w, "데이터베이스 오류", http.StatusInternalServerError)
		return
	}
	log.Printf("[GetLeaveStatus] 사용자 입사일 문자열 조회 성공: %v", hireDateStr)
	log.Printf("[GetLeaveStatus] 입사일 문자열 타입: %T", hireDateStr)

	// 문자열을 time.Time으로 변환
	hireDate, err := time.Parse("2006-01-02", hireDateStr)
	if err != nil {
		log.Printf("[GetLeaveStatus] 입사일 파싱 오류: %v", err)
		http.Error(w, "입사일 형식 오류", http.StatusInternalServerError)
		return
	}
	log.Printf("[GetLeaveStatus] 입사일 변환 성공: %v", hireDate)
	log.Printf("[GetLeaveStatus] 변환된 입사일 타입: %T", hireDate)

	// 사용한 연차 일수 조회
	var usedDays int
	log.Printf("[GetLeaveStatus] 연차 사용 내역 조회 시작: userID=%d, year=%d", userID, year)
	log.Printf("[GetLeaveStatus] 연차 사용 내역 쿼리 실행")

	err = config.DB.QueryRow(`
		SELECT COALESCE(SUM(DATEDIFF(end_date, start_date) + 1), 0) 
		FROM leave_requests 
		WHERE user_id = ? AND YEAR(start_date) = ? AND leave_type = '연차' AND status = '승인'
	`, userID, year).Scan(&usedDays)
	if err != nil {
		log.Printf("[GetLeaveStatus] 연차 사용 내역 조회 오류: %v", err)
		log.Printf("[GetLeaveStatus] 오류 타입: %T", err)
		http.Error(w, "연차 사용 내역 조회 오류", http.StatusInternalServerError)
		return
	}
	log.Printf("[GetLeaveStatus] 사용한 연차 일수 조회 성공: %d", usedDays)
	log.Printf("[GetLeaveStatus] usedDays 타입: %T", usedDays)

	// 연차 계산기 생성
	lc := utils.NewLeaveCalculator()
	log.Printf("[GetLeaveStatus] 연차 계산기 생성 완료")

	// 연차 현황 계산
	log.Printf("[GetLeaveStatus] 연차 계산 시작: userID=%d, hireDate=%v, usedDays=%d", userID, hireDate, usedDays)
	leaveStatus := lc.CalculateUserLeaveStatus(userID, hireDate, usedDays)
	log.Printf("[GetLeaveStatus] 연차 현황 계산 완료: total=%d, used=%d, remaining=%d",
		leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)
	log.Printf("[GetLeaveStatus] leaveStatus 타입 확인: TotalEarned=%T, UsedDays=%T, RemainingDays=%T",
		leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)

	// DB에 연차 현황 저장/업데이트
	log.Printf("[GetLeaveStatus] DB 저장 시작")
	log.Printf("[GetLeaveStatus] 저장할 데이터: userID=%d, totalEarned=%v, usedDays=%v, remainingDays=%v",
		leaveStatus.UserID, leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)
	log.Printf("[GetLeaveStatus] 데이터 타입: userID=%T, totalEarned=%T, usedDays=%T, remainingDays=%T",
		leaveStatus.UserID, leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)

	_, err = config.DB.Exec(`
		INSERT INTO user_leave_status (user_id, total_earned_days, used_days, remaining_days)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			total_earned_days = VALUES(total_earned_days),
			used_days = VALUES(used_days),
			remaining_days = VALUES(remaining_days),
			updated_at = CURRENT_TIMESTAMP
	`, leaveStatus.UserID, leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)

	if err != nil {
		log.Printf("[GetLeaveStatus] DB 저장 오류: %v", err)
		log.Printf("[GetLeaveStatus] 오류 타입: %T", err)
		log.Printf("[GetLeaveStatus] 오류 상세: %+v", err)
		http.Error(w, "연차 현황 저장 오류", http.StatusInternalServerError)
		return
	}
	log.Printf("[GetLeaveStatus] DB 저장 완료")

	// JSON 응답
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    leaveStatus,
	})
}

// GetLeavePolicy 연차 정책 정보 조회
func GetLeavePolicy(w http.ResponseWriter, r *http.Request) {
	lc := utils.NewLeaveCalculator()
	policy := lc.GetLeavePolicy()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    policy,
	})
}

// CalculateLeaveForDate 특정 날짜 기준 연차 계산
func CalculateLeaveForDate(w http.ResponseWriter, r *http.Request) {
	// URL 파라미터에서 user_id와 target_date 추출
	userIDStr := r.URL.Query().Get("user_id")
	targetDateStr := r.URL.Query().Get("target_date")

	if userIDStr == "" || targetDateStr == "" {
		http.Error(w, "user_id와 target_date 파라미터가 필요합니다", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "잘못된 user_id입니다", http.StatusBadRequest)
		return
	}

	targetDate, err := time.Parse("2006-01-02", targetDateStr)
	if err != nil {
		http.Error(w, "잘못된 target_date 형식입니다 (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// 사용자 정보 조회 (입사일 포함)
	var hireDateStr string
	log.Printf("[CalculateLeaveForDate] 사용자 정보 조회 시작: userID=%d", userID)
	log.Printf("[CalculateLeaveForDate] DB 쿼리 실행: SELECT hire_date FROM users WHERE id = %d", userID)

	err = config.DB.QueryRow("SELECT hire_date FROM users WHERE id = ?", userID).Scan(&hireDateStr)
	if err != nil {
		log.Printf("[CalculateLeaveForDate] 사용자 정보 조회 오류: %v", err)
		log.Printf("[CalculateLeaveForDate] 오류 타입: %T", err)
		if err == sql.ErrNoRows {
			log.Printf("[CalculateLeaveForDate] 사용자를 찾을 수 없음: userID=%d", userID)
			http.Error(w, "사용자를 찾을 수 없습니다", http.StatusNotFound)
			return
		}
		log.Printf("[CalculateLeaveForDate] 데이터베이스 오류 발생")
		http.Error(w, "데이터베이스 오류", http.StatusInternalServerError)
		return
	}
	log.Printf("[CalculateLeaveForDate] 사용자 입사일 문자열 조회 성공: %v", hireDateStr)
	log.Printf("[CalculateLeaveForDate] 입사일 문자열 타입: %T", hireDateStr)

	// 문자열을 time.Time으로 변환
	hireDate, err := time.Parse("2006-01-02", hireDateStr)
	if err != nil {
		log.Printf("[CalculateLeaveForDate] 입사일 파싱 오류: %v", err)
		http.Error(w, "입사일 형식 오류", http.StatusInternalServerError)
		return
	}
	log.Printf("[CalculateLeaveForDate] 입사일 변환 성공: %v", hireDate)
	log.Printf("[CalculateLeaveForDate] 변환된 입사일 타입: %T", hireDate)

	// 연차 계산
	lc := utils.NewLeaveCalculator()
	log.Printf("[CalculateLeaveForDate] 연차 계산기 생성 완료")

	log.Printf("[CalculateLeaveForDate] 연차 계산 시작: hireDate=%v, targetDate=%v", hireDate, targetDate)
	yearsOfService := lc.CalculateYearsOfService(hireDate, targetDate)
	log.Printf("[CalculateLeaveForDate] 근무 연수 계산 완료: %d년", yearsOfService)

	totalLeave := lc.CalculateTotalLeave(hireDate, targetDate)
	log.Printf("[CalculateLeaveForDate] 총 연차 계산 완료: %d일", totalLeave)

	// 응답 데이터
	result := map[string]interface{}{
		"user_id":          userID,
		"hire_date":        hireDate.Format("2006-01-02"),
		"target_date":      targetDate.Format("2006-01-02"),
		"years_of_service": yearsOfService,
		"total_leave_days": totalLeave,
		"calculation_type": "under_one_year",
	}

	if yearsOfService >= 1 {
		result["calculation_type"] = "over_one_year"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// UpdateAllUsersLeaveStatus 모든 사용자의 연차 현황 업데이트
func UpdateAllUsersLeaveStatus(w http.ResponseWriter, r *http.Request) {
	// 현재 연도 사용
	year := time.Now().Year()

	// 모든 사용자 조회
	log.Printf("[UpdateAllUsersLeaveStatus] 모든 사용자 조회 시작")
	rows, err := config.DB.Query("SELECT id, hire_date FROM users")
	if err != nil {
		log.Printf("[UpdateAllUsersLeaveStatus] 사용자 목록 조회 오류: %v", err)
		log.Printf("[UpdateAllUsersLeaveStatus] 오류 타입: %T", err)
		http.Error(w, "사용자 목록 조회 오류", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	log.Printf("[UpdateAllUsersLeaveStatus] 사용자 목록 조회 성공")

	lc := utils.NewLeaveCalculator()
	updatedCount := 0
	for rows.Next() {
		var userID int
		var hireDateStr string

		err := rows.Scan(&userID, &hireDateStr)
		if err != nil {
			log.Printf("[UpdateAllUsersLeaveStatus] 사용자 데이터 스캔 오류: %v", err)
			continue
		}
		log.Printf("[UpdateAllUsersLeaveStatus] 사용자 처리 시작: userID=%d, hireDateStr=%v", userID, hireDateStr)

		// 문자열을 time.Time으로 변환
		hireDate, err := time.Parse("2006-01-02", hireDateStr)
		if err != nil {
			log.Printf("[UpdateAllUsersLeaveStatus] 입사일 파싱 오류: userID=%d, error=%v", userID, err)
			continue
		}
		log.Printf("[UpdateAllUsersLeaveStatus] 입사일 변환 성공: userID=%d, hireDate=%v", userID, hireDate)

		// 사용한 연차 일수 조회
		var usedDays int
		log.Printf("[UpdateAllUsersLeaveStatus] 연차 사용 내역 조회: userID=%d, year=%d", userID, year)
		err = config.DB.QueryRow(`
			SELECT COALESCE(SUM(DATEDIFF(end_date, start_date) + 1), 0) 
			FROM leave_requests 
			WHERE user_id = ? AND YEAR(start_date) = ? AND leave_type = '연차' AND status = '승인'
		`, userID, year).Scan(&usedDays)
		if err != nil {
			log.Printf("[UpdateAllUsersLeaveStatus] 연차 사용 내역 조회 오류: userID=%d, error=%v", userID, err)
			continue
		}
		log.Printf("[UpdateAllUsersLeaveStatus] 연차 사용 내역 조회 성공: userID=%d, usedDays=%d", userID, usedDays)

		// 연차 현황 계산
		log.Printf("[UpdateAllUsersLeaveStatus] 연차 현황 계산 시작: userID=%d, hireDate=%v, usedDays=%d", userID, hireDate, usedDays)
		leaveStatus := lc.CalculateUserLeaveStatus(userID, hireDate, usedDays)
		log.Printf("[UpdateAllUsersLeaveStatus] 연차 현황 계산 완료: userID=%d, total=%d, used=%d, remaining=%d",
			userID, leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)

		// DB에 저장/업데이트
		log.Printf("[UpdateAllUsersLeaveStatus] DB 저장 시작: userID=%d", userID)
		_, err = config.DB.Exec(`
			INSERT INTO user_leave_status (user_id, total_earned_days, used_days, remaining_days)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				total_earned_days = VALUES(total_earned_days),
				used_days = VALUES(used_days),
				remaining_days = VALUES(remaining_days),
				updated_at = CURRENT_TIMESTAMP
		`, leaveStatus.UserID, leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)

		if err == nil {
			log.Printf("[UpdateAllUsersLeaveStatus] DB 저장 성공: userID=%d", userID)
			updatedCount++
		} else {
			log.Printf("[UpdateAllUsersLeaveStatus] DB 저장 오류: userID=%d, error=%v", userID, err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"message":       "연차 현황 업데이트 완료",
		"updated_count": updatedCount,
		"year":          year,
	})
}
