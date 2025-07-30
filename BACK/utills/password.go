package utills

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// GenerateSalt는 임의의 솔트를 생성하는 함수입니다
func GenerateSalt(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// HashPasswordWithSalt는 비밀번호와 솔트를 결합하여 해싱하는 함수입니다
func HashPasswordWithSalt(password string, salt string) (string, error) {
	// 솔트와 비밀번호를 합쳐서 해싱
	saltedPassword := password + salt
	bytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), 14)
	return string(bytes), err
}

// CheckPasswordWithSalt는 솔트를 사용하여 비밀번호와 해시를 비교하는 함수입니다
func CheckPasswordWithSalt(password string, salt string, hash string) bool {
	saltedPassword := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltedPassword))
	return err == nil
}

//공통 유틸리티 함수 모음  , 프로젝트 전역에서 자주 쓰는 기능들을 한데 모음
