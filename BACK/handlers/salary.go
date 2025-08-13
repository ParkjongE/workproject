package handlers

import (
	"BACK/models"
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
	GetSalaryStatement(userID int, year int, month int) (*models.SalaryFile, error)
	GetSalaryStatementList(userID int) ([]models.SalaryStatement, error)
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
	fmt.Printf("[핸들러] GetSalaryStatementList 호출됨! URL: %s\n", r.URL.String())
	userIDStr := r.URL.Query().Get("user_id")
	fmt.Printf("[핸들러] userIDStr: %s\n", userIDStr)
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
func (s *SalaryServiceImpl) GetSalaryStatement(userID int, year int, month int) (*models.SalaryFile, error) {
	fileName := fmt.Sprintf("salary_statement_%d_%d_%d.pdf", userID, year, month)
	filePath := filepath.Join(s.basePath, fileName)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &models.SalaryFile{
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    fileInfo.Size(),
		ContentType: "application/pdf",
	}, nil
}

// GetSalaryStatementList 서비스 구현
func (s *SalaryServiceImpl) GetSalaryStatementList(userID int) ([]models.SalaryStatement, error) {
	var statements []models.SalaryStatement

	fmt.Printf("[GetSalaryStatementList] 파일 검색 시작: userID=%d, basePath=%s\n", userID, s.basePath)

	files, err := os.ReadDir(s.basePath)
	if err != nil {
		fmt.Printf("[GetSalaryStatementList] 디렉토리 읽기 오류: %v\n", err)
		return nil, err
	}

	fmt.Printf("[GetSalaryStatementList] 총 파일 개수: %d\n", len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		fmt.Printf("[GetSalaryStatementList] 파일 확인: %s\n", fileName)

		if strings.HasPrefix(fileName, fmt.Sprintf("salary_statement_%d_", userID)) {
			fmt.Printf("[GetSalaryStatementList] 매칭된 파일: %s\n", fileName)
			parts := strings.Split(strings.TrimSuffix(fileName, ".pdf"), "_")
			fmt.Printf("[GetSalaryStatementList] 파일명 분리: %v\n", parts)

			if len(parts) == 5 {
				year, _ := strconv.Atoi(parts[3])
				month, _ := strconv.Atoi(parts[4])

				statement := models.SalaryStatement{
					UserID:   userID,
					Year:     year,
					Month:    month,
					FileName: fileName,
					FilePath: filepath.Join(s.basePath, fileName),
				}
				fmt.Printf("[GetSalaryStatementList] 추가된 statement: %+v\n", statement)
				statements = append(statements, statement)
			}
		}
	}

	fmt.Printf("[GetSalaryStatementList] 최종 결과: %d개 파일\n", len(statements))
	return statements, nil
}
