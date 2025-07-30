package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"BACK/config"
	"BACK/models"
	"BACK/utills"
)

// JWTLoginRequest는 JWT 로그인 요청 구조체입니다
type JWTLoginRequest struct {
	Userid   string `json:"userid"`
	Password string `json:"password"`
}

// JWTLoginResponse는 JWT 로그인 응답 구조체입니다
type JWTLoginResponse struct {
	Success bool        `json:"success"`
	Token   string      `json:"token"`
	User    models.User `json:"user"`
	Message string      `json:"message"`
}

// HandleJWTLogin은 JWT 기반 로그인을 처리하는 핸들러 함수입니다
func HandleJWTLogin(w http.ResponseWriter, r *http.Request) {
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

	// 데이터베이스에서 사용자 조회
	var user models.User
	var storedPassword string
	var salt string

	err = config.DB.QueryRow("SELECT id, userid, password, salt, name, role, department FROM users WHERE userid = ?",
		loginData.Userid).Scan(&user.ID, &user.Userid, &storedPassword, &salt, &user.Name, &user.Role, &user.Department)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("사용자 없음: %s", loginData.Userid)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "사용자명 또는 비밀번호가 일치하지 않습니다"})
			return
		}
		log.Println("사용자 조회 오류:", err)
		http.Error(w, "서버 오류", http.StatusInternalServerError)
		return
	}

	// 비밀번호 검증
	passwordValid := utills.CheckPasswordWithSalt(loginData.Password, salt, storedPassword)

	if !passwordValid {
		log.Println("비밀번호 불일치")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "사용자명 또는 비밀번호가 일치하지 않습니다"})
		return
	}

	log.Printf("JWT 로그인 성공: %s (%s)", user.Userid, user.Name)

	// JWT 토큰 생성
	token, err := utills.GenerateJWT(user.ID, user.Userid, user.Role)
	if err != nil {
		log.Printf("JWT 토큰 생성 실패: %v", err)
		http.Error(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}

	// 응답 생성
	response := JWTLoginResponse{
		Success: true,
		Token:   token,
		User:    user,
		Message: "JWT 로그인 성공",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleJWTLogout은 JWT 로그아웃을 처리하는 핸들러 함수입니다
// JWT는 서버 측에서 무효화할 수 없으므로, 클라이언트에서 토큰을 삭제하도록 안내합니다
func HandleJWTLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "로그아웃 성공. 클라이언트에서 토큰을 삭제하세요.",
	})
}

// HandleTokenRefresh는 JWT 토큰 갱신을 처리하는 핸들러 함수입니다
func HandleTokenRefresh(w http.ResponseWriter, r *http.Request) {
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

	tokenString, err := utills.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Printf("토큰 추출 실패: %v", err)
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	// 토큰 갱신
	newToken, err := utills.RefreshJWT(tokenString)
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
func HandleVerifyToken(w http.ResponseWriter, r *http.Request) {
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

	tokenString, err := utills.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Printf("토큰 추출 실패: %v", err)
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	// 토큰 검증
	claims, err := utills.ValidateJWT(tokenString)
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
