package model

import (
	"BACK/utils" //
	"database/sql"
	"net/http"
	"time"
)

// AuthService 인증 관련 서비스 인터페이스
type AuthService interface {
	Login(userid, password string) (*User, string, string, error)
	Logout(sessionID string) error
	RefreshToken(refreshToken string) (string, string, error)
	GetCurrentUser(userID int) (*User, error)
	ValidatePassword(password, salt, storedPassword string) bool
	GetUserProfile(w http.ResponseWriter, r *http.Request) // 이 줄 추가
}

// UserDB DB에서 직접 읽어오는 사용자 구조체
type UserDB struct {
	ID         int          `db:"id"`
	Userid     string       `db:"userid"`
	Password   string       `db:"password"`
	Salt       string       `db:"salt"`
	Name       string       `db:"name"`
	Role       string       `db:"role"`
	Department string       `db:"department"`
	CreatedAt  sql.NullTime `db:"created_at"`
	HireDate   time.Time    `db:"hire_date"`
	// UpdatedAt 필드 제거 (현재 테이블에 없음)
}

// ToUser DB 구조체를 API 응답 구조체로 변환
func (db *UserDB) ToUser() *User {
	var createdAt *time.Time

	if db.CreatedAt.Valid {
		createdAt = &db.CreatedAt.Time
	}

	return &User{
		ID:         db.ID,
		Userid:     db.Userid,
		Name:       db.Name,
		Role:       db.Role,
		Department: db.Department,
		HireDate:   db.HireDate,
		CreatedAt:  createdAt,
		UpdatedAt:  nil, // UpdatedAt은 현재 테이블에 없음
	}
}

// UserRegistrationRequest 사용자 등록 요청 모델
type UserRegistrationRequest struct {
	Userid     string    `json:"userid"`
	Password   string    `json:"password"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Department string    `json:"department"`
	HireDate   time.Time `json:"hire_date"`
}

// UserRegistrationResponse 사용자 등록 응답 모델
type UserRegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ValidatePassword(password, salt, storedPassword string) bool {
	return utils.VerifyPassword(password, salt, storedPassword)
}
