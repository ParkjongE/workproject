package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetSalaryStatement 급여명세서 조회
func GetSalaryStatement(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GetSalaryStatement] 요청 시작")

	// URL 파라미터에서 user_id, year, month 추출
	userIDStr := r.URL.Query().Get("user_id")
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	log.Printf("[GetSalaryStatement] 파라미터: user_id=%s, year=%s, month=%s", userIDStr, yearStr, monthStr)

	if userIDStr == "" || yearStr == "" || monthStr == "" {
		log.Printf("[GetSalaryStatement] 필수 파라미터 누락")
		http.Error(w, "user_id, year, month 파라미터가 필요합니다", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("[GetSalaryStatement] user_id 변환 오류: %v", err)
		http.Error(w, "잘못된 user_id입니다", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Printf("[GetSalaryStatement] year 변환 오류: %v", err)
		http.Error(w, "잘못된 year입니다", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		log.Printf("[GetSalaryStatement] month 변환 오류: %v", err)
		http.Error(w, "잘못된 month입니다", http.StatusBadRequest)
		return
	}

	// 파일 경로 생성
	fileName := fmt.Sprintf("salary_statement_%d_%d_%d.pdf", userID, year, month)
	filePath := filepath.Join("files", "salary", fileName)

	log.Printf("[GetSalaryStatement] 파일 경로: %s", filePath)

	// 파일 존재 여부 확인
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("[GetSalaryStatement] 파일이 존재하지 않음: %s", filePath)
		http.Error(w, "급여명세서를 찾을 수 없습니다", http.StatusNotFound)
		return
	}

	// 파일 읽기
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("[GetSalaryStatement] 파일 열기 오류: %v", err)
		http.Error(w, "파일을 열 수 없습니다", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 파일 정보 가져오기
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("[GetSalaryStatement] 파일 정보 조회 오류: %v", err)
		http.Error(w, "파일 정보를 가져올 수 없습니다", http.StatusInternalServerError)
		return
	}

	// HTTP 헤더 설정
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	log.Printf("[GetSalaryStatement] 파일 전송 시작: %s, 크기: %d bytes", fileName, fileInfo.Size())

	// 파일 내용을 응답에 복사
	http.ServeContent(w, r, fileName, fileInfo.ModTime(), file)

	log.Printf("[GetSalaryStatement] 파일 전송 완료")
}

// GetSalaryStatementList 급여명세서 목록 조회
func GetSalaryStatementList(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GetSalaryStatementList] 요청 시작")
	log.Printf("[GetSalaryStatementList] 요청 URL: %s", r.URL.String())
	log.Printf("[GetSalaryStatementList] 요청 메서드: %s", r.Method)
	log.Printf("[GetSalaryStatementList] 요청 헤더: %v", r.Header)

	// URL 파라미터에서 user_id 추출
	userIDStr := r.URL.Query().Get("user_id")

	log.Printf("[GetSalaryStatementList] 파라미터: user_id=%s", userIDStr)

	if userIDStr == "" {
		log.Printf("[GetSalaryStatementList] user_id 파라미터 누락")
		http.Error(w, "user_id 파라미터가 필요합니다", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("[GetSalaryStatementList] user_id 변환 오류: %v", err)
		http.Error(w, "잘못된 user_id입니다", http.StatusBadRequest)
		return
	}

	// 급여명세서 디렉토리 경로
	salaryDir := filepath.Join("files", "salary")
	log.Printf("[GetSalaryStatementList] 급여명세서 디렉토리 경로: %s", salaryDir)

	// 디렉토리 존재 여부 확인
	if _, err := os.Stat(salaryDir); os.IsNotExist(err) {
		log.Printf("[GetSalaryStatementList] 급여명세서 디렉토리가 존재하지 않음: %s", salaryDir)
		// 빈 목록 반환
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "data": []}`))
		return
	}
	log.Printf("[GetSalaryStatementList] 급여명세서 디렉토리 존재 확인됨")

	// 디렉토리 읽기
	log.Printf("[GetSalaryStatementList] 디렉토리 읽기 시작")
	entries, err := os.ReadDir(salaryDir)
	if err != nil {
		log.Printf("[GetSalaryStatementList] 디렉토리 읽기 오류: %v", err)
		http.Error(w, "디렉토리를 읽을 수 없습니다", http.StatusInternalServerError)
		return
	}
	log.Printf("[GetSalaryStatementList] 디렉토리 읽기 성공, 파일 수: %d", len(entries))

	// 모든 파일명 로그
	for i, entry := range entries {
		log.Printf("[GetSalaryStatementList] 파일 %d: %s (디렉토리: %v)", i, entry.Name(), entry.IsDir())
	}

	// 해당 사용자의 급여명세서 파일들 필터링
	var statements []map[string]interface{}
	prefix := fmt.Sprintf("salary_statement_%d_", userID)
	log.Printf("[GetSalaryStatementList] 검색할 파일명 접두사: %s", prefix)

	for _, entry := range entries {
		log.Printf("[GetSalaryStatementList] 파일 검사: %s", entry.Name())
		log.Printf("[GetSalaryStatementList] - 디렉토리 여부: %v", entry.IsDir())
		log.Printf("[GetSalaryStatementList] - 접두사 일치: %v", strings.HasPrefix(entry.Name(), prefix))
		log.Printf("[GetSalaryStatementList] - PDF 확장자: %v", strings.HasSuffix(entry.Name(), ".pdf"))

		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) && strings.HasSuffix(entry.Name(), ".pdf") {
			// 파일명에서 년도와 월 추출 (salary_statement_1_2024_3.pdf)
			parts := strings.Split(strings.TrimSuffix(entry.Name(), ".pdf"), "_")
			if len(parts) >= 5 {
				year, _ := strconv.Atoi(parts[3])  // parts[3] = 2024
				month, _ := strconv.Atoi(parts[4]) // parts[4] = 3

				// 파일 정보 가져오기
				filePath := filepath.Join(salaryDir, entry.Name())
				fileInfo, err := os.Stat(filePath)
				fileSize := int64(0)
				if err == nil {
					fileSize = fileInfo.Size()
				}

				statements = append(statements, map[string]interface{}{
					"year":     year,
					"month":    month,
					"fileName": entry.Name(),
					"fileSize": fileSize,
				})
			}
		}
	}

	log.Printf("[GetSalaryStatementList] 조회된 급여명세서 수: %d", len(statements))
	log.Printf("[GetSalaryStatementList] 최종 statements: %+v", statements)

	// JSON 응답
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"data":    statements,
	}

	log.Printf("[GetSalaryStatementList] 응답 전송: %+v", response)
	json.NewEncoder(w).Encode(response)
	log.Printf("[GetSalaryStatementList] 응답 전송 완료")
}
