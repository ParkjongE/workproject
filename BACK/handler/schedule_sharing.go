package handler

import (
	"BACK/model"
	"BACK/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// ScheduleSharingHandler 스케줄 공유 관련 HTTP 요청 처리
type ScheduleSharingHandler struct {
	service service.ScheduleSharingService
}

// NewScheduleSharingHandler ScheduleSharingHandler 생성자
func NewScheduleSharingHandler(service service.ScheduleSharingService) *ScheduleSharingHandler {
	return &ScheduleSharingHandler{service: service}
}

// GetUsersByDepartment 부서별 사용자 조회 핸들러
func (h *ScheduleSharingHandler) GetUsersByDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 부서명 추출: /api/schedule-sharing/users/department/개발부
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid department path", http.StatusBadRequest)
		return
	}

	department := pathParts[4]

	// 부서명 유효성 검사
	validDepartments := map[string]bool{
		"개발부":  true,
		"디자인부": true,
		"마케팅부": true,
	}

	if !validDepartments[department] {
		http.Error(w, "Invalid department", http.StatusBadRequest)
		return
	}

	// 서비스를 통해 부서별 사용자 정보 조회
	users, err := h.service.GetUsersByDepartment(department)
	if err != nil {
		response := model.ScheduleSharingResponse{
			Success: false,
			Error:   fmt.Sprintf("부서별 사용자 조회 실패: %v", err),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 성공 응답
	response := model.ScheduleSharingResponse{
		Success: true,
		Message: fmt.Sprintf("%s 부서 사용자 조회 성공", department),
		Data:    users,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetAllDepartmentUsers 모든 부서 사용자 조회 핸들러
func (h *ScheduleSharingHandler) GetAllDepartmentUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 서비스를 통해 모든 부서 사용자 정보 조회
	users, err := h.service.GetAllDepartmentUsers()
	if err != nil {
		response := model.ScheduleSharingResponse{
			Success: false,
			Error:   fmt.Sprintf("전체 부서 사용자 조회 실패: %v", err),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 성공 응답
	response := model.ScheduleSharingResponse{
		Success: true,
		Message: "전체 부서 사용자 조회 성공",
		Data:    users,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetUserSchedulesByID 특정 사용자의 일정만 조회
func (h *ScheduleSharingHandler) GetUserSchedulesByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	schedules, err := h.service.GetUserSchedulesByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// 일정이 없는 경우 빈 배열 반환
			response := map[string]interface{}{
				"success": true,
				"data":    []model.Schedule{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data":    schedules,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserWithSchedulesByID 특정 사용자의 일정 조회 핸들러
func (h *ScheduleSharingHandler) GetUserWithSchedulesByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 사용자 ID 추출: /api/schedule-sharing/user/{userID}/schedules
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response := model.ScheduleResponse{
			Success: false,
			Error:   "Invalid user ID",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 서비스를 통해 사용자 일정 조회
	schedules, err := h.service.GetUserSchedulesByID(userID)
	if err != nil {
		response := model.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("사용자 %d 일정 조회 실패: %v", userID, err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 성공 응답
	response := model.ScheduleResponse{
		Success: true,
		Message: fmt.Sprintf("사용자 %d 일정 조회 성공", userID),
		Data:    schedules,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
