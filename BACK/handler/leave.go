package handler

import (
	"BACK/model"
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

// LeaveHandler 연차 핸들러
type LeaveHandler struct {
	service LeaveService
}

// NewLeaveHandler 연차 핸들러 생성자
func NewLeaveHandler(service LeaveService) *LeaveHandler {
	return &LeaveHandler{service: service}
}

// GetUserLeaveStatus 사용자 연차 현황 조회 (순수 HTTP 처리)
func (h *LeaveHandler) GetUserLeaveStatus(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetUserLeaveStatus(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.APIResponse{
		Success: true,
		Data:    status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CalculateLeaveForDate 특정 날짜 연차 계산 (순수 HTTP 처리)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.LeaveCalculationResponse{
		Success: true,
		Data:    result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateAllUsersLeaveStatus 모든 사용자 연차 현황 업데이트 (순수 HTTP 처리)
func (h *LeaveHandler) UpdateAllUsersLeaveStatus(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.UpdateAllUsersLeaveStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := model.LeaveUpdateResponse{
		Success: true,
		Message: "Updated leave status for all users",
		Count:   count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
