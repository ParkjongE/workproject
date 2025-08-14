package service

import (
	"BACK/model"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// SalaryServiceImpl 급여 서비스 구현체
type SalaryServiceImpl struct {
	basePath string
}

// NewSalaryService 급여 서비스 생성자
func NewSalaryService(basePath string) *SalaryServiceImpl {
	return &SalaryServiceImpl{basePath: basePath}
}

// GetSalaryStatement 급여명세서 파일 조회 (비즈니스 로직)
func (s *SalaryServiceImpl) GetSalaryStatement(userID, year, month int) (*model.SalaryFile, error) {
	// 비즈니스 검증
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}
	if year < 2020 || year > 2030 {
		return nil, fmt.Errorf("invalid year")
	}
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month")
	}

	fileName := fmt.Sprintf("salary_statement_%d_%d_%d.pdf", userID, year, month)
	filePath := filepath.Join(s.basePath, fileName)

	// 파일 존재 확인
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("salary statement not found")
	}

	return &model.SalaryFile{
		FileName:    fileName,
		FilePath:    filePath,
		ContentType: "application/pdf",
	}, nil
}

// GetSalaryStatementList 급여명세서 목록 조회 (비즈니스 로직)
func (s *SalaryServiceImpl) GetSalaryStatementList(userID int) ([]model.SalaryStatement, error) {
	// 비즈니스 검증
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	var statements []model.SalaryStatement

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

				statement := model.SalaryStatement{
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
