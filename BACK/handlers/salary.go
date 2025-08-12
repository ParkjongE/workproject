package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// SalaryService 급여 관련 서비스 인터페이스
type SalaryService interface {
	GetSalaryStatement(userID int, year int, month int) (*SalaryFile, error)
	GetSalaryStatementList(userID int) ([]SalaryStatement, error)
}

// SalaryServiceImpl 급여 서비스 구현체
type SalaryServiceImpl struct {
	basePath string
}

// NewSalaryService 급여 서비스 생성자
func NewSalaryService(basePath string) SalaryService {
	return &SalaryServiceImpl{basePath: basePath}
}

// SalaryHandler 급여 핸들러
type SalaryHandler struct {
	service SalaryService
}

// NewSalaryHandler 급여 핸들러 생성자
func NewSalaryHandler(service SalaryService) *SalaryHandler {
	return &SalaryHandler{service: service}
}

// SalaryFile 급여 파일 정보 구조체
type SalaryFile struct {
	FileName    string
	FilePath    string
	FileSize    int64
	ContentType string
}

// SalaryStatement 급여 명세서 정보 구조체
type SalaryStatement struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}

// GetSalaryStatement 급여 명세서 조회
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
		http.Error(w, "Failed to get salary statement", http.StatusInternalServerError)
		return
	}

	// 파일 다운로드
	http.ServeFile(w, r, file.FilePath)
}

// GetSalaryStatementList 급여 명세서 목록 조회
func (h *SalaryHandler) GetSalaryStatementList(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	list, err := h.service.GetSalaryStatementList(userID)
	if err != nil {
		http.Error(w, "Failed to get salary statement list", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data":    list,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSalaryStatement 서비스 구현
func (s *SalaryServiceImpl) GetSalaryStatement(userID int, year int, month int) (*SalaryFile, error) {
	fileName := fmt.Sprintf("salary_statement_%d_%d_%d.pdf", userID, year, month)
	filePath := filepath.Join(s.basePath, fileName)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &SalaryFile{
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    fileInfo.Size(),
		ContentType: "application/pdf",
	}, nil
}

// GetSalaryStatementList 서비스 구현
func (s *SalaryServiceImpl) GetSalaryStatementList(userID int) ([]SalaryStatement, error) {
	var statements []SalaryStatement

	files, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if strings.HasPrefix(fileName, fmt.Sprintf("salary_statement_%d_", userID)) {
			parts := strings.Split(strings.TrimSuffix(fileName, ".pdf"), "_")
			if len(parts) == 4 {
				year, _ := strconv.Atoi(parts[2])
				month, _ := strconv.Atoi(parts[3])

				statements = append(statements, SalaryStatement{
					UserID:   userID,
					Year:     year,
					Month:    month,
					FileName: fileName,
					FilePath: filepath.Join(s.basePath, fileName),
				})
			}
		}
	}

	return statements, nil
}
