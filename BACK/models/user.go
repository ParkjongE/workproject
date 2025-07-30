package models

import "time"

// User는 시스템 사용자 정보를 나타내는 구조체
type User struct {
	ID         int       `json:"id"`
	Userid     string    `json:"userid"`
	Password   string    `json:"password"` // 요청에서 받을 수 있도록 json 태그 수정
	Salt       string    `json:"salt"`     // JSON 응답에서 제외 (솔트)
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Department string    `json:"department"`
	HireDate   time.Time `json:"hire_date"`
}

// LoginRequest는 로그인 요청 데이터를 나타내는 구조체
type LoginRequest struct {
	Userid   string `json:"userid"`
	Password string `json:"password"`
}
