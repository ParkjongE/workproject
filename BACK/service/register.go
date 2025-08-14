package service

import (
	"BACK/model"
	"BACK/repository"
	"BACK/utils"
	"time"
)

// RegisterError 사용자 등록 관련 에러 구조체
type RegisterError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *RegisterError) Error() string {
	return e.Message
}

// RegisterServiceImpl 사용자 등록 서비스 구현체
type RegisterServiceImpl struct {
	repo repository.RegisterRepository
}

// NewRegisterService 사용자 등록 서비스 생성자
func NewRegisterService(repo repository.RegisterRepository) *RegisterServiceImpl {
	return &RegisterServiceImpl{repo: repo}
}

// ValidateUserData 사용자 데이터 검증 (비즈니스 로직)
func (s *RegisterServiceImpl) ValidateUserData(user *model.User) error {
	if user.Userid == "" || user.Password == "" || user.Name == "" {
		return &RegisterError{Message: "사용자명, 비밀번호, 이름은 필수 입력 항목입니다", Code: "MISSING_REQUIRED_FIELDS"}
	}

	// 비즈니스 규칙: userid 길이 검증
	if len(user.Userid) < 3 {
		return &RegisterError{Message: "사용자명은 3자 이상이어야 합니다", Code: "USERID_TOO_SHORT"}
	}

	// 비즈니스 규칙: password 길이 검증
	if len(user.Password) < 4 {
		return &RegisterError{Message: "비밀번호는 4자 이상이어야 합니다", Code: "PASSWORD_TOO_SHORT"}
	}

	// hire_date가 비어있으면 기본값 설정 (2024-01-01)
	if user.HireDate.IsZero() {
		user.HireDate = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	// role이 비어있으면 기본값 "사원"으로 설정
	if user.Role == "" {
		user.Role = "사원"
	}

	return nil
}

// CheckUserExists 사용자명 중복 확인 (비즈니스 로직)
func (s *RegisterServiceImpl) CheckUserExists(userid string) error {
	existingUser, err := s.repo.FindUserByUserid(userid)
	if err != nil {
		return &RegisterError{Message: "서버 오류", Code: "DB_ERROR"}
	}

	if existingUser != nil {
		return &RegisterError{Message: "이미 사용 중인 사용자명입니다", Code: "USER_ALREADY_EXISTS"}
	}

	return nil
}

// HashUserPassword 사용자 비밀번호 해싱 (비즈니스 로직)
func (s *RegisterServiceImpl) HashUserPassword(password string) (string, string, error) {
	// 1. 솔트 생성
	salt := utils.GenerateSalt()

	// 2. 비밀번호와 솔트를 결합하여 해싱
	hashedPassword := utils.HashPassword(password, salt)

	return hashedPassword, salt, nil
}

// RegisterUser 사용자 등록 (비즈니스 로직)
func (s *RegisterServiceImpl) RegisterUser(user *model.User) error {
	// hire_date가 비어있으면 기본값 설정 (2024-01-01)
	if user.HireDate.IsZero() {
		user.HireDate = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	// Repository를 통해 사용자 생성
	err := s.repo.CreateUser(user)
	if err != nil {
		return &RegisterError{Message: "사용자 등록에 실패했습니다", Code: "DB_ERROR"}
	}

	return nil
}
