package utils

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// JWT 클레임 구조체
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Userid string `json:"userid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWT 시크릿 키 (환경변수에서 가져오거나 기본값 사용)
var jwtSecret []byte

// InitJWTSecret은 JWT 시크릿 키를 초기화합니다
func InitJWTSecret() {
	// 현재 실행 경로 확인
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("현재 디렉토리 확인 실패:", err)
	} else {
		log.Printf("현재 작업 디렉토리: %s", currentDir)
	}

	// 여러 경로에서 .env 파일 로드 시도
	envPaths := []string{
		filepath.Join(currentDir, ".env"),
		filepath.Join(currentDir, "BACK", ".env"),
		filepath.Join(currentDir, "..", "BACK", ".env"),
		filepath.Join(currentDir, "..", "..", "BACK", ".env"),
		filepath.Join("C:", "project", ".env"),
		filepath.Join("C:", "project", "BACK", ".env"),
	}

	envLoaded := false
	for _, path := range envPaths {
		log.Printf(".env 파일 시도: %s", path)
		if _, err := os.Stat(path); err == nil {
			log.Printf(".env 파일 발견: %s", path)
			if err := godotenv.Load(path); err == nil {
				log.Printf("JWT: .env 파일을 로드했습니다: %s", path)
				envLoaded = true
				break
			} else {
				log.Printf(".env 파일 로드 실패: %s - %v", path, err)
			}
		}
	}

	if !envLoaded {
		log.Println("경고: JWT 설정을 위한 .env 파일을 찾을 수 없습니다. 시스템 환경변수를 사용합니다.")
	}

	// 환경변수에서 JWT 시크릿 키 가져오기 (여러 키 이름 시도)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET_KEY")
	}
	if secret == "" {
		secret = "your-super-secret-jwt-key-change-this-in-production"
		log.Println("경고: JWT_SECRET 또는 JWT_SECRET_KEY 환경변수가 설정되지 않았습니다. 기본값을 사용합니다.")
	} else {
		log.Println("JWT 시크릿 키를 환경변수에서 로드했습니다.")
	}

	jwtSecret = []byte(secret)
}

// GenerateJWT는 사용자 정보로 JWT 토큰을 생성합니다
func GenerateJWT(userID int, userid, role string) (string, error) {
	// JWT 시크릿이 초기화되지 않았다면 초기화
	if jwtSecret == nil {
		InitJWTSecret()
	}

	// 토큰 만료 시간 설정 (환경변수에서 가져오거나 기본값 24시간)
	expirationHours := 24
	if envHours := os.Getenv("JWT_EXPIRATION_HOURS"); envHours != "" {
		if hours, err := strconv.Atoi(envHours); err == nil && hours > 0 {
			expirationHours = hours
		}
	}
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// 클레임 생성
	claims := &JWTClaims{
		UserID: userID,
		Userid: userid,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "project-api",
			Subject:   userid,
		},
	}

	// 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateAccessToken은 Access Token을 생성합니다 (15분 만료)
func GenerateAccessToken(userID int, userid, role string) (string, error) {
	log.Printf("GenerateAccessToken 호출: userID=%d, userid=%s, role=%s", userID, userid, role)

	// JWT 시크릿이 초기화되지 않았다면 초기화
	if jwtSecret == nil {
		log.Printf("JWT 시크릿이 nil, 초기화 시작")
		InitJWTSecret()
	}

	log.Printf("JWT 시크릿 길이: %d", len(jwtSecret))

	// Access Token 만료 시간 설정 (15분)
	expirationTime := time.Now().Add(15 * time.Minute)

	// 클레임 생성
	claims := &JWTClaims{
		UserID: userID,
		Userid: userid,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "project-api",
			Subject:   userid,
		},
	}

	// 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken은 Refresh Token을 생성합니다 (7일 만료)
func GenerateRefreshToken(userID int, userid string) (string, error) {
	// JWT 시크릿이 초기화되지 않았다면 초기화
	if jwtSecret == nil {
		InitJWTSecret()
	}

	// Refresh Token 만료 시간 설정 (7일)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Refresh Token용 클레임 (role 제외)
	claims := &JWTClaims{
		UserID: userID,
		Userid: userid,
		Role:   "", // Refresh Token에는 role 정보 제외
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "project-api",
			Subject:   userid,
		},
	}

	// 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT는 JWT 토큰을 검증하고 클레임을 반환합니다
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	// JWT 시크릿이 초기화되지 않았다면 초기화
	if jwtSecret == nil {
		InitJWTSecret()
	}

	// 토큰 파싱
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 서명 방법 검증
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 클레임 추출
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshAccessToken은 Refresh Token으로 새로운 Access Token을 생성합니다
func RefreshAccessToken(refreshTokenString string) (string, error) {
	// Refresh Token 검증
	claims, err := ValidateJWT(refreshTokenString)
	if err != nil {
		return "", err
	}

	// Refresh Token으로 새로운 Access Token 생성
	// role 정보는 데이터베이스에서 다시 조회하거나,
	// Refresh Token에 포함된 정보를 사용
	accessToken, err := GenerateAccessToken(claims.UserID, claims.Userid, claims.Role)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// GetUserInfoFromRefreshToken은 Refresh Token에서 사용자 정보를 추출합니다
func GetUserInfoFromRefreshToken(refreshTokenString string) (int, string, error) {
	// Refresh Token 검증
	claims, err := ValidateJWT(refreshTokenString)
	if err != nil {
		return 0, "", err
	}

	return claims.UserID, claims.Userid, nil
}

// ExtractTokenFromHeader는 Authorization 헤더에서 JWT 토큰을 추출합니다
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}
	return authHeader[7:], nil
}

// RefreshJWT는 기존 토큰을 검증하고 새로운 토큰을 생성합니다
func RefreshJWT(tokenString string) (string, error) {
	// 기존 토큰 검증
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}

	// 새로운 토큰 생성
	newToken, err := GenerateJWT(claims.UserID, claims.Userid, claims.Role)
	if err != nil {
		return "", err
	}

	return newToken, nil
}
