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
	var hireDate time.Time
	log.Printf("[GetLeaveStatus] 사용자 정보 조회 시작: userID=%d", userID)
	err = config.DB.QueryRow("SELECT hire_date FROM users WHERE id = ?", userID).Scan(&hireDate)
	if err != nil {
		log.Printf("[GetLeaveStatus] 사용자 정보 조회 오류: %v", err)
		if err == sql.ErrNoRows {
			http.Error(w, "사용자를 찾을 수 없습니다", http.StatusNotFound)
			return
		}
		http.Error(w, "데이터베이스 오류", http.StatusInternalServerError)
		return
	}
	log.Printf("[GetLeaveStatus] 사용자 입사일: %v", hireDate)

	// 사용한 연차 일수 조회
	var usedDays int
	log.Printf("[GetLeaveStatus] 연차 사용 내역 조회 시작: userID=%d, year=%d", userID, year)
	err = config.DB.QueryRow(`
		SELECT COALESCE(SUM(DATEDIFF(end_date, start_date) + 1), 0) 
		FROM leave_requests 
		WHERE user_id = ? AND YEAR(start_date) = ? AND leave_type = '연차' AND status = '승인'
	`, userID, year).Scan(&usedDays)
	if err != nil {
		log.Printf("[GetLeaveStatus] 연차 사용 내역 조회 오류: %v", err)
		http.Error(w, "연차 사용 내역 조회 오류", http.StatusInternalServerError)
		return
	}
	log.Printf("[GetLeaveStatus] 사용한 연차 일수: %d", usedDays)

	// 연차 계산기 생성
	lc := utils.NewLeaveCalculator()
	log.Printf("[GetLeaveStatus] 연차 계산기 생성 완료")

	// 연차 현황 계산
	leaveStatus := lc.CalculateUserLeaveStatus(userID, hireDate, usedDays)
	log.Printf("[GetLeaveStatus] 연차 현황 계산 완료: total=%d, used=%d, remaining=%d",
		leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)

	// DB에 연차 현황 저장/업데이트
	log.Printf("[GetLeaveStatus] DB 저장 시작")
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
	var hireDate time.Time
	err = config.DB.QueryRow("SELECT hire_date FROM users WHERE id = ?", userID).Scan(&hireDate)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "사용자를 찾을 수 없습니다", http.StatusNotFound)
			return
		}
		http.Error(w, "데이터베이스 오류", http.StatusInternalServerError)
		return
	}

	// 연차 계산
	lc := utils.NewLeaveCalculator()
	yearsOfService := lc.CalculateYearsOfService(hireDate, targetDate)
	totalLeave := lc.CalculateTotalLeave(hireDate, targetDate)

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
	rows, err := config.DB.Query("SELECT id, hire_date FROM users")
	if err != nil {
		http.Error(w, "사용자 목록 조회 오류", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	lc := utils.NewLeaveCalculator()
	updatedCount := 0
	for rows.Next() {
		var userID int
		var hireDate time.Time

		err := rows.Scan(&userID, &hireDate)
		if err != nil {
			continue
		}

		// 사용한 연차 일수 조회
		var usedDays int
		err = config.DB.QueryRow(`
			SELECT COALESCE(SUM(DATEDIFF(end_date, start_date) + 1), 0) 
			FROM leave_requests 
			WHERE user_id = ? AND YEAR(start_date) = ? AND leave_type = '연차' AND status = '승인'
		`, userID, year).Scan(&usedDays)
		if err != nil {
			continue
		}

		// 연차 현황 계산
		leaveStatus := lc.CalculateUserLeaveStatus(userID, hireDate, usedDays)

		// DB에 저장/업데이트
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
			updatedCount++
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
