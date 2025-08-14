package model

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

// ==================== Response 구조체들 ====================

// SalaryListResponse 급여명세서 목록 응답 구조체
type SalaryListResponse struct {
	Success bool              `json:"success"`
	Data    []SalaryStatement `json:"data"`
}

// SalaryStatementResponse 급여명세서 파일 응답 구조체
type SalaryStatementResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
}

// SalaryErrorResponse 급여 관련 에러 응답 구조체
type SalaryErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
