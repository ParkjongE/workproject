package handlers

import (
	"BACK/config"
	"BACK/models"
	"encoding/json"
	"net/http"
)

// ScheduleService 스케줄 관련 서비스 인터페이스
type ScheduleService interface {
	CreateSchedule(schedule *models.Schedule) error
	GetSchedulesByUser(userID string) ([]models.Schedule, error)
	UpdateSchedule(scheduleID string, schedule *models.Schedule) error
	DeleteSchedule(scheduleID string) error
}

// ScheduleServiceImpl 스케줄 서비스 구현체
type ScheduleServiceImpl struct {
	db config.Database
}

// NewScheduleService 스케줄 서비스 생성자
func NewScheduleService(db config.Database) ScheduleService {
	return &ScheduleServiceImpl{db: db}
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
	var schedule models.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.CreateSchedule(&schedule)
	if err != nil {
		http.Error(w, "Failed to create schedule", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Schedule created successfully",
		"data":    schedule,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSchedulesByUser 사용자별 스케줄 조회
func (h *ScheduleHandler) GetSchedulesByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	schedules, err := h.service.GetSchedulesByUser(userID)
	if err != nil {
		http.Error(w, "Failed to get schedules", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}

// UpdateSchedule 스케줄 수정
func (h *ScheduleHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID := r.URL.Query().Get("schedule_id")
	if scheduleID == "" {
		http.Error(w, "Schedule ID is required", http.StatusBadRequest)
		return
	}

	var schedule models.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateSchedule(scheduleID, &schedule)
	if err != nil {
		http.Error(w, "Failed to update schedule", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Schedule updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteSchedule 스케줄 삭제
func (h *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID := r.URL.Query().Get("schedule_id")
	if scheduleID == "" {
		http.Error(w, "Schedule ID is required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteSchedule(scheduleID)
	if err != nil {
		http.Error(w, "Failed to delete schedule", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Schedule deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateSchedule 서비스 구현
func (s *ScheduleServiceImpl) CreateSchedule(schedule *models.Schedule) error {
	query := "INSERT INTO schedules (user_id, title, content, start_time, end_time, color, created_at) VALUES (?, ?, ?, ?, ?, ?, NOW())"

	_, err := s.db.Exec(query, schedule.UserID, schedule.Title, schedule.Content, schedule.StartTime, schedule.EndTime, schedule.Color)
	return err
}

// GetSchedulesByUser 서비스 구현
func (s *ScheduleServiceImpl) GetSchedulesByUser(userID string) ([]models.Schedule, error) {
	query := "SELECT schedule_id, user_id, title, content, start_time, end_time, created_at, color FROM schedules WHERE user_id = ? ORDER BY start_time"

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.ScheduleID, &schedule.UserID, &schedule.Title, &schedule.Content,
			&schedule.StartTime, &schedule.EndTime, &schedule.CreatedAt, &schedule.Color,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// UpdateSchedule 서비스 구현
func (s *ScheduleServiceImpl) UpdateSchedule(scheduleID string, schedule *models.Schedule) error {
	query := "UPDATE schedules SET title = ?, content = ?, start_time = ?, end_time = ?, color = ? WHERE schedule_id = ?"

	_, err := s.db.Exec(query, schedule.Title, schedule.Content, schedule.StartTime, schedule.EndTime, schedule.Color, scheduleID)
	return err
}

// DeleteSchedule 서비스 구현
func (s *ScheduleServiceImpl) DeleteSchedule(scheduleID string) error {
	query := "DELETE FROM schedules WHERE schedule_id = ?"

	_, err := s.db.Exec(query, scheduleID)
	return err
}
