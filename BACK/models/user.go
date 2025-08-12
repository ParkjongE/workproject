package models

import (
	"time"
)

// User API 응답용 사용자 구조체
type User struct {
	ID         int        `json:"id"`
	Userid     string     `json:"userid"`
	Password   string     `json:"-"`
	Salt       string     `json:"-"`
	Name       string     `json:"name"`
	Role       string     `json:"role"`
	Department string     `json:"department"`
	HireDate   time.Time  `json:"hire_date"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// LoginRequest 로그인 요청 모델 (auth.go에서 제거하고 여기로 이동)
type LoginRequest struct {
	Userid   string `json:"userid"`
	Password string `json:"password"`
}

// LoginResponse 로그인 응답 모델 (auth.go에서 제거하고 여기로 이동)
type LoginResponse struct {
	Success      bool   `json:"success"`
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}

// RefreshTokenRequest 토큰 갱신 요청 모델 (auth.go에서 제거하고 여기로 이동)
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse 토큰 갱신 응답 모델 (auth.go에서 제거하고 여기로 이동)
type RefreshTokenResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}
