package models

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
