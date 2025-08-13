package handlers

import (
	"BACK/config"
	"BACK/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// LeaveService 연차 관련 서비스 인터페이스
type LeaveService interface {
	GetUserLeaveStatus(userID int) (*utils.LeaveStatus, error)
	GetLeavePolicy() *utils.LeavePolicy
	CalculateLeaveForDate(userID int, targetDate time.Time) (map[string]interface{}, error)
	UpdateAllUsersLeaveStatus() (int, error)
}

// LeaveServiceImpl 연차 서비스 구현체
type LeaveServiceImpl struct {
	db config.Database
	lc *utils.LeaveCalculator
}

// NewLeaveService 연차 서비스 생성자
func NewLeaveService(db config.Database) LeaveService {
	return &LeaveServiceImpl{
		db: db,
		lc: utils.NewLeaveCalculator(),
	}
}

// LeaveHandler 연차 핸들러
type LeaveHandler struct {
	service LeaveService
}

// NewLeaveHandler 연차 핸들러 생성자
func NewLeaveHandler(service LeaveService) *LeaveHandler {
	return &LeaveHandler{service: service}
}

// GetUserLeaveStatus 사용자 연차 현황 조회
func (h *LeaveHandler) GetUserLeaveStatus(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetUserLeaveStatus(userID)
	if err != nil {
		http.Error(w, "Failed to get leave status", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data":    status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CalculateLeaveForDate 특정 날짜 연차 계산
func (h *LeaveHandler) CalculateLeaveForDate(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	result, err := h.service.CalculateLeaveForDate(userID, date)
	if err != nil {
		http.Error(w, "Failed to calculate leave", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data":    result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateAllUsersLeaveStatus 모든 사용자 연차 현황 업데이트
func (h *LeaveHandler) UpdateAllUsersLeaveStatus(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.UpdateAllUsersLeaveStatus()
	if err != nil {
		http.Error(w, "Failed to update leave status", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Updated leave status for all users",
		"count":   count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserLeaveStatus 서비스 구현
func (s *LeaveServiceImpl) GetUserLeaveStatus(userID int) (*utils.LeaveStatus, error) {
	// 1. 사용자 정보 가져오기 (hire_date 포함)
	var hireDate time.Time
	var usedDays int

	userQuery := "SELECT hire_date FROM users WHERE id = ?"
	err := s.db.QueryRow(userQuery, userID).Scan(&hireDate)
	if err != nil {
		return nil, err
	}

	// 2. 사용된 연차 계산 (승인된 연차 신청 기준)
	usedQuery := `
		SELECT COALESCE(SUM(DATEDIFF(end_date, start_date) + 1), 0) 
		FROM leave_requests 
		WHERE user_id = ? AND leave_type = '연차' AND status = '승인' 
		AND YEAR(start_date) = YEAR(CURDATE())
	`
	err = s.db.QueryRow(usedQuery, userID).Scan(&usedDays)
	if err != nil {
		usedDays = 0 // 사용된 연차가 없으면 0
	}

	// 3. 연차 계산기로 총 연차 계산
	leaveStatus := s.lc.CalculateUserLeaveStatus(userID, hireDate, usedDays)

	// 4. DB에 현재 연차 상태 업데이트/삽입
	updateQuery := `
		INSERT INTO user_leave_status (user_id, total_earned_days, used_days, remaining_days, updated_at)
		VALUES (?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
		total_earned_days = VALUES(total_earned_days),
		used_days = VALUES(used_days),
		remaining_days = VALUES(remaining_days),
		updated_at = NOW()
	`

	_, err = s.db.Exec(updateQuery, userID, leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)
	if err != nil {
		return nil, err
	}

	return leaveStatus, nil
}

// GetLeavePolicy 서비스 구현
func (s *LeaveServiceImpl) GetLeavePolicy() *utils.LeavePolicy {
	return s.lc.GetLeavePolicy()
}

// CalculateLeaveForDate 서비스 구현
func (s *LeaveServiceImpl) CalculateLeaveForDate(userID int, targetDate time.Time) (map[string]interface{}, error) {
	// 연차 계산 로직 구현
	result := map[string]interface{}{
		"user_id":     userID,
		"target_date": targetDate.Format("2006-01-02"),
		"is_holiday":  false,
		"leave_type":  "none",
	}

	return result, nil
}

// UpdateAllUsersLeaveStatus 서비스 구현
func (s *LeaveServiceImpl) UpdateAllUsersLeaveStatus() (int, error) {
	// 모든 사용자 연차 현황 업데이트 로직 구현
	query := "UPDATE leave_status SET remaining_leave = total_leave - used_leave WHERE year = YEAR(CURDATE())"

	result, err := s.db.Exec(query)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
