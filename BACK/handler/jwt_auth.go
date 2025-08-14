package handler

import (
	"BACK/model"
	"BACK/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// JWTService JWT 관련 서비스 인터페이스
type JWTService interface {
	AuthenticateUser(userid, password string) (*model.User, error)
	GenerateToken(user *model.User) (string, error)
	RefreshToken(tokenString string) (string, error)
	ValidateToken(tokenString string) (*utils.JWTClaims, error)
}

// JWTServiceImpl JWT 서비스 구현체
type JWTServiceImpl struct {
	db *sql.DB
}

// NewJWTService JWT 서비스 생성자
func NewJWTService(db *sql.DB) JWTService {
	return &JWTServiceImpl{db: db}
}

// AuthenticateUser 사용자 인증 처리
func (s *JWTServiceImpl) AuthenticateUser(userid, password string) (*model.User, error) {
	var user model.User
	var storedPassword string
	var salt string

	err := s.db.QueryRow("SELECT id, userid, password, salt, name, role, department FROM users WHERE userid = ?",
		userid).Scan(&user.ID, &user.Userid, &storedPassword, &salt, &user.Name, &user.Role, &user.Department)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &AuthError{Message: "사용자명 또는 비밀번호가 일치하지 않습니다", Code: "USER_NOT_FOUND"}
		}
		return nil, &AuthError{Message: "사용자 조회 오류", Code: "DB_ERROR"}
	}

	// 비밀번호 검증
	passwordValid := utils.VerifyPassword(password, salt, storedPassword)
	if !passwordValid {
		return nil, &AuthError{Message: "사용자명 또는 비밀번호가 일치하지 않습니다", Code: "INVALID_PASSWORD"}
	}

	return &user, nil
}

// GenerateToken JWT 토큰 생성
func (s *JWTServiceImpl) GenerateToken(user *model.User) (string, error) {
	return utils.GenerateJWT(user.ID, user.Userid, user.Role)
}

// RefreshToken JWT 토큰 갱신
func (s *JWTServiceImpl) RefreshToken(tokenString string) (string, error) {
	return utils.RefreshJWT(tokenString)
}

// ValidateToken JWT 토큰 검증
func (s *JWTServiceImpl) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	return utils.ValidateJWT(tokenString)
}

// JWTLoginRequest는 JWT 로그인 요청 구조체입니다
type JWTLoginRequest struct {
	Userid   string `json:"userid"`
	Password string `json:"password"`
}

// JWTLoginResponse는 JWT 로그인 응답 구조체입니다
type JWTLoginResponse struct {
	Success bool       `json:"success"`
	Token   string     `json:"token"`
	User    model.User `json:"user"`
	Message string     `json:"message"`
}

// AuthError 인증 관련 에러 구조체
type AuthError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *AuthError) Error() string {
	return e.Message
}

// JWTHandler JWT 핸들러 구조체
type JWTHandler struct {
	service JWTService
}

// NewJWTHandler JWT 핸들러 생성자
func NewJWTHandler(service JWTService) *JWTHandler {
	return &JWTHandler{service: service}
}

// HandleJWTLogin은 JWT 기반 로그인을 처리하는 핸들러 함수입니다
func (h *JWTHandler) HandleJWTLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("잘못된 HTTP 메서드: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 요청 본문에서 로그인 정보 파싱
	var loginData JWTLoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		log.Println("JWT 로그인 요청 파싱 오류:", err)
		http.Error(w, "잘못된 요청 형식", http.StatusBadRequest)
		return
	}

	log.Printf("JWT 로그인 시도: 사원번호=%s", loginData.Userid)

	// 사용자 인증
	user, err := h.service.AuthenticateUser(loginData.Userid, loginData.Password)
	if err != nil {
		if authErr, ok := err.(*AuthError); ok {
			log.Printf("인증 실패: %s", authErr.Message)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": authErr.Message})
			return
		}
		log.Println("서버 오류:", err)
		http.Error(w, "서버 오류", http.StatusInternalServerError)
		return
	}

	log.Printf("JWT 로그인 성공: %s (%s)", user.Userid, user.Name)

	// JWT 토큰 생성
	token, err := h.service.GenerateToken(user)
	if err != nil {
		log.Printf("JWT 토큰 생성 실패: %v", err)
		http.Error(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}

	// 응답 생성
	response := JWTLoginResponse{
		Success: true,
		Token:   token,
		User:    *user,
		Message: "JWT 로그인 성공",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleJWTLogout은 JWT 로그아웃을 처리하는 핸들러 함수입니다
func (h *JWTHandler) HandleJWTLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "로그아웃 성공. 클라이언트에서 토큰을 삭제하세요.",
	})
}

// HandleTokenRefresh는 JWT 토큰 갱신을 처리하는 핸들러 함수입니다
func (h *JWTHandler) HandleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("잘못된 HTTP 메서드: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authorization 헤더에서 토큰 추출
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization 헤더가 없습니다")
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Printf("토큰 추출 실패: %v", err)
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	// 토큰 갱신
	newToken, err := h.service.RefreshToken(tokenString)
	if err != nil {
		log.Printf("토큰 갱신 실패: %v", err)
		http.Error(w, "Token refresh failed", http.StatusUnauthorized)
		return
	}

	// 응답 생성
	response := map[string]interface{}{
		"success": true,
		"token":   newToken,
		"message": "토큰 갱신 성공",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleVerifyToken은 JWT 토큰 검증을 처리하는 핸들러 함수입니다
func (h *JWTHandler) HandleVerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("잘못된 HTTP 메서드: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authorization 헤더에서 토큰 추출
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization 헤더가 없습니다")
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Printf("토큰 추출 실패: %v", err)
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	// 토큰 검증
	claims, err := h.service.ValidateToken(tokenString)
	if err != nil {
		log.Printf("토큰 검증 실패: %v", err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// 응답 생성
	response := map[string]interface{}{
		"success": true,
		"valid":   true,
		"user": map[string]interface{}{
			"user_id": claims.UserID,
			"userid":  claims.Userid,
			"role":    claims.Role,
		},
		"message": "토큰이 유효합니다",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
