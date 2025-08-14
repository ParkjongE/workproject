package handler

import (
	"BACK/config"
	"BACK/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"BACK/utils"
)

// AuthService 인증 관련 서비스 인터페이스
type AuthService interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	ValidateUser(userid, password string) (*model.User, error)
	RefreshToken(tokenString string) (string, error)
	GetCurrentUser(tokenString string) (*model.User, error)
	GetUserProfile(w http.ResponseWriter, r *http.Request) // 이 줄 추가
}

// AuthServiceImpl 인증 서비스 구현체
type AuthServiceImpl struct {
	db config.Database
}

// NewAuthService 인증 서비스 생성자
func NewAuthService(db config.Database) AuthService {
	return &AuthServiceImpl{db: db}
}

// AuthHandler 인증 핸들러
type AuthHandler struct {
	service AuthService
}

// NewAuthHandler 인증 핸들러 생성자
func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// HandleLogin 로그인 처리
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	h.service.HandleLogin(w, r)
}

// ValidateUser 사용자 검증
func (s *AuthServiceImpl) ValidateUser(userid, password string) (*model.User, error) {
	var userDB model.UserDB
	query := "SELECT id, userid, password, salt, name, role, department, created_at, hire_date FROM users WHERE userid = ?"

	err := s.db.QueryRow(query, userid).Scan(
		&userDB.ID, &userDB.Userid, &userDB.Password, &userDB.Salt,
		&userDB.Name, &userDB.Role, &userDB.Department,
		&userDB.CreatedAt, &userDB.HireDate,
	)

	if err != nil {
		return nil, err
	}

	// 비밀번호 검증
	if !utils.VerifyPassword(password, userDB.Salt, userDB.Password) {
		return nil, sql.ErrNoRows
	}

	return userDB.ToUser(), nil
}

// RefreshToken 토큰 갱신
func (s *AuthServiceImpl) RefreshToken(tokenString string) (string, error) {
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}

	// 새 토큰 생성
	newToken, err := utils.GenerateJWT(claims.UserID, claims.Userid, claims.Role)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

// GetCurrentUser 현재 사용자 정보 조회
func (s *AuthServiceImpl) GetCurrentUser(tokenString string) (*model.User, error) {
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	var userDB model.UserDB
	query := "SELECT id, userid, password, salt, name, role, department, created_at, upda	ted_at FROM users WHERE id = ?"

	err = s.db.QueryRow(query, claims.UserID).Scan(
		&userDB.ID, &userDB.Userid, &userDB.Password, &userDB.Salt,
		&userDB.Name, &userDB.Role, &userDB.Department,
		&userDB.CreatedAt, &userDB.HireDate,
	)

	if err != nil {
		return nil, err
	}

	return userDB.ToUser(), nil
}

// HandleLogin 로그인 처리 (구현체에 직접 구현)
func (s *AuthServiceImpl) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var err error // 함수 시작 부분에 err 변수 선언

	var loginReq model.LoginRequest
	if err = json.NewDecoder(r.Body).Decode(&loginReq); err != nil { // := 대신 = 사용
		log.Printf("JSON 파싱 오류: %v", err)
		log.Printf("요청 헤더: %v", r.Header)
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("파싱된 로그인 요청: userid=%s, password=[숨김]", loginReq.Userid)

	// 사용자 조회
	var userDB model.UserDB
	query := `SELECT id, userid, password, salt, name, role, department, hire_date, created_at FROM users WHERE userid = ?`
	err = s.db.QueryRow(query, loginReq.Userid).Scan( // := 대신 = 사용
		&userDB.ID, &userDB.Userid, &userDB.Password, &userDB.Salt,
		&userDB.Name, &userDB.Role, &userDB.Department, &userDB.CreatedAt,
		&userDB.HireDate,
	)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// JWT 토큰 생성
	token, err := utils.GenerateJWT(userDB.ID, userDB.Userid, userDB.Role) // 이건 := 유지 (새로운 변수)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// JWT 토큰 생성 후 응답
	response := map[string]interface{}{
		"success":      true,
		"access_token": token, // token → access_token으로 변경
		"user":         userDB.ToUser(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshTokenHandler 토큰 갱신 핸들러
func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// h.service를 통해 메서드 호출 (구현체 인스턴스 사용)
	newToken, err := h.service.RefreshToken(req.Token)
	if err != nil {
		http.Error(w, "Token refresh failed", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"token":   newToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUserHandler 현재 사용자 정보 조회 핸들러
func (h *AuthHandler) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Authorization 헤더에서 토큰 추출
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	// "Bearer " 제거
	tokenString := authHeader[7:]

	// h.service를 통해 메서드 호출 (구현체 인스턴스 사용)
	user, err := h.service.GetCurrentUser(tokenString)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"user":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleLogout 로그아웃 핸들러
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, "user-session")
	if err != nil {
		log.Printf("세션 가져오기 실패: %v", err)
		http.Error(w, "세션 오류", http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = false
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		log.Printf("세션 저장 오류: %v", err)
		http.Error(w, "서버 오류", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "로그아웃 성공",
	})
}

// HandleRefreshToken 토큰 갱신 핸들러
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("잘못된 HTTP 메서드: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var refreshToken string
	var requestBody model.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err == nil {
		refreshToken = requestBody.RefreshToken
	}

	if refreshToken == "" {
		refreshTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			log.Printf("Refresh Token 없음: %v", err)
			http.Error(w, "Refresh token not found", http.StatusUnauthorized)
			return
		}
		refreshToken = refreshTokenCookie.Value
	}

	// JWT 토큰 검증 및 갱신
	claims, err := utils.ValidateJWT(refreshToken)
	if err != nil {
		log.Printf("Access Token 갱신 실패: %v", err)
		http.Error(w, "Token refresh failed", http.StatusUnauthorized)
		return
	}

	// 새 토큰 생성
	newAccessToken, err := utils.GenerateJWT(claims.UserID, claims.Userid, claims.Role)
	if err != nil {
		log.Printf("새 토큰 생성 실패: %v", err)
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	response := model.RefreshTokenResponse{
		Success:      true,
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken, // 기존 Refresh Token 사용
		Message:      "Access Token 갱신 성공",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser 현재 사용자 정보 조회 핸들러
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "세션 오류", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "로그인이 필요합니다", http.StatusUnauthorized)
		return
	}

	// 데이터베이스에서 사용자 정보 조회
	var userDB model.UserDB
	query := "SELECT id, userid, password, salt, name, role, department, created_at, updated_at FROM users WHERE id = ?"

	err = config.DB.QueryRow(query, userID).Scan(
		&userDB.ID, &userDB.Userid, &userDB.Password, &userDB.Salt,
		&userDB.Name, &userDB.Role, &userDB.Department,
		&userDB.CreatedAt, &userDB.HireDate,
	)

	if err != nil {
		log.Printf("오류: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := userDB.ToUser()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// GetUserProfile 사용자 프로필 조회
func (s *AuthServiceImpl) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("=== GetUserProfile 요청 시작 ===")
	log.Printf("요청 메서드: %s", r.Method)
	log.Printf("요청 URL: %s", r.URL.String())

	// 요청 헤더 로깅
	log.Println("요청 헤더:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("  %s: %s", name, value)
		}
	}

	if r.Method != http.MethodGet {
		log.Println("잘못된 HTTP 메서드")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// JWT 토큰에서 사용자 ID 추출
	log.Println("JWT 토큰에서 사용자 ID 추출 시작")
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		log.Printf("토큰에서 사용자 ID 추출 실패: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("사용자 ID 추출 성공: %d", userID)

	// 사용자 정보 조회
	log.Println("사용자 정보 조회 시작")
	var userDB model.UserDB
	query := `SELECT id, userid, password, salt, name, role, department, hire_date, created_at FROM users WHERE id = ?`
	err = s.db.QueryRow(query, userID).Scan(
		&userDB.ID, &userDB.Userid, &userDB.Password, &userDB.Salt,
		&userDB.Name, &userDB.Role, &userDB.Department, &userDB.CreatedAt,
		&userDB.HireDate,
	)

	if err != nil {
		log.Printf("사용자 조회 실패: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	log.Printf("사용자 조회 성공: %s", userDB.Name)

	// 응답
	log.Println("응답 생성 및 전송")
	response := map[string]interface{}{
		"success": true,
		"user":    userDB.ToUser(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Println("=== GetUserProfile 요청 완료 ===")
}
