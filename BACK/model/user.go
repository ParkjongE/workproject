package model

import (
	"time"
)

// User 사용자 정보 구조체 (DB용)
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

// RegisterRequest 회원가입 요청 구조체 (JSON용)
type RegisterRequest struct {
	Userid     string `json:"userid"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Department string `json:"department"`
	HireDate   string `json:"hire_date"` // 문자열로 받기
}

// ToUser RegisterRequest를 User로 변환
func (req *RegisterRequest) ToUser() (*User, error) {
	// 날짜 파싱
	hireDate, err := time.Parse("2006-01-02", req.HireDate)
	if err != nil {
		return nil, err
	}

	return &User{
		Userid:     req.Userid,
		Password:   req.Password,
		Name:       req.Name,
		Role:       req.Role,
		Department: req.Department,
		HireDate:   hireDate,
	}, nil
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

// ==================== 공통 Response 구조체들 ====================

// APIResponse 공통 API 응답 구조체
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 성공 응답 구조체
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse 에러 응답 구조체
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
