package handler

import (
	"BACK/model"
	"encoding/json"
	"net/http"
)

// ScheduleService 스케줄 관련 서비스 인터페이스
type ScheduleService interface {
	CreateSchedule(schedule *model.Schedule) error
	GetSchedulesByUser(userID string) ([]model.Schedule, error)
	UpdateSchedule(scheduleID string, schedule *model.Schedule) error
	DeleteSchedule(scheduleID string) error
}

// ScheduleHandler 스케줄 핸들러
type ScheduleHandler struct {
	service ScheduleService
}

// NewScheduleHandler 스케줄 핸들러 생성자
func NewScheduleHandler(service ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

// CreateSchedule 스케줄 생성
func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var schedule model.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.CreateSchedule(&schedule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.ScheduleAPIResponse{
		Success: true,
		Message: "Schedule created successfully",
		Data:    &schedule,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSchedulesByUser 사용자별 스케줄 조회
func (h *ScheduleHandler) GetSchedulesByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	schedules, err := h.service.GetSchedulesByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}

// UpdateSchedule 스케줄 수정
func (h *ScheduleHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID := r.URL.Query().Get("schedule_id")

	var schedule model.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateSchedule(scheduleID, &schedule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.ScheduleSuccessResponse{
		Success: true,
		Message: "Schedule updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteSchedule 스케줄 삭제
func (h *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID := r.URL.Query().Get("schedule_id")

	err := h.service.DeleteSchedule(scheduleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.ScheduleSuccessResponse{
		Success: true,
		Message: "Schedule deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
