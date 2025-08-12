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
	var status utils.LeaveStatus
	query := "SELECT user_id, total_earned, used_days, remaining_days FROM leave_status WHERE user_id = ? AND year = YEAR(CURDATE())"

	err := s.db.QueryRow(query, userID).Scan(
		&status.UserID, &status.TotalEarned, &status.UsedDays, &status.RemainingDays,
	)

	if err != nil {
		return nil, err
	}

	return &status, nil
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
