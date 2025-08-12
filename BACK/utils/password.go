package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateSalt 랜덤 salt 생성 (32바이트)
func GenerateSalt() string {
	salt := make([]byte, 32)
	rand.Read(salt)
	return hex.EncodeToString(salt)
}

// HashPassword 비밀번호에 salt를 추가하여 SHA-256으로 해싱
func HashPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

// VerifyPassword 비밀번호 검증
func VerifyPassword(password, salt, hashedPassword string) bool {
	expectedHash := HashPassword(password, salt)
	return expectedHash == hashedPassword
}
