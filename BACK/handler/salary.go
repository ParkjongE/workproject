package handler

import (
	"BACK/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// SalaryService 급여 관련 서비스 인터페이스
type SalaryService interface {
	GetSalaryStatement(userID, year, month int) (*model.SalaryFile, error)
	GetSalaryStatementList(userID int) ([]model.SalaryStatement, error)
}

// SalaryHandler 급여 핸들러
type SalaryHandler struct {
	service SalaryService
}

// NewSalaryHandler 급여 핸들러 생성자
func NewSalaryHandler(service SalaryService) *SalaryHandler {
	return &SalaryHandler{service: service}
}

// GetSalaryStatement 급여명세서 파일 다운로드
func (h *SalaryHandler) GetSalaryStatement(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	file, err := h.service.GetSalaryStatement(userID, year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))
	http.ServeFile(w, r, file.FilePath)
}

// GetSalaryStatementList 급여명세서 목록 조회
func (h *SalaryHandler) GetSalaryStatementList(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	list, err := h.service.GetSalaryStatementList(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.SalaryListResponse{
		Success: true,
		Data:    list,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
