package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"BACK/config"
	"BACK/middleware"
	"BACK/models"
	"BACK/utills"
)

// HandleLogin은 로그인 요청을 처리하는 핸들러 함수입니다
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("잘못된 HTTP 메서드: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 요청 본문에서 로그인 정보 파싱
	var loginData models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		log.Println("로그인 요청 파싱 오류:", err)
		http.Error(w, "잘못된 요청 형식", http.StatusBadRequest)
		return
	}

	log.Printf("로그인 시도: 사원번호=%s", loginData.Userid)

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
			json.NewEncoder(w).Encode(map[string]string{"error": "사원번호 또는 비밀번호가 일치하지 않습니다"})
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
		json.NewEncoder(w).Encode(map[string]string{"error": "사원번호 또는 비밀번호가 일치하지 않습니다"})
		return
	}

	log.Printf("로그인 성공: %s (%s)", user.Userid, user.Name)
	log.Printf("사용자 정보: ID=%d, Userid=%s, Role=%s", user.ID, user.Userid, user.Role)

	// 세션 생성
	session, err := config.Store.Get(r, "user-session")
	if err != nil {
		log.Printf("세션 생성 오류: %v", err)
		http.Error(w, "서버 오류", http.StatusInternalServerError)
		return
	}

	// 세션에 사용자 정보 저장
	session.Values["user_id"] = user.ID
	session.Values["userid"] = user.Userid
	session.Values["username"] = user.Userid // username 필드 추가 (DualAuthMiddleware와 일치)
	session.Values["role"] = user.Role
	session.Values["authenticated"] = true

	// 세션 저장
	err = session.Save(r, w)
	if err != nil {
		log.Printf("세션 저장 오류: %v", err)
		http.Error(w, "서버 오류", http.StatusInternalServerError)
		return
	}

	log.Printf("세션이 성공적으로 저장되었습니다. 사원번호: %s", user.Userid)

	// Access Token 생성
	log.Printf("Access Token 생성 시작: userID=%d, userid=%s, role=%s", user.ID, user.Userid, user.Role)
	accessToken, err := utills.GenerateAccessToken(user.ID, user.Userid, user.Role)
	if err != nil {
		log.Printf("Access Token 생성 실패: %v", err)
		http.Error(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}
	log.Printf("Access Token 생성 성공: %s...", accessToken[:20])

	// Refresh Token 생성
	log.Printf("Refresh Token 생성 시작: userID=%d, userid=%s", user.ID, user.Userid)
	refreshToken, err := utills.GenerateRefreshToken(user.ID, user.Userid)
	if err != nil {
		log.Printf("Refresh Token 생성 실패: %v", err)
		http.Error(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}
	log.Printf("Refresh Token 생성 성공: %s...", refreshToken[:20])

	// Refresh Token을 HttpOnly 쿠키로 설정
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   7 * 24 * 3600, // 7일
		HttpOnly: true,
		Secure:   false, // HTTPS 환경에서는 true로 설정
		SameSite: http.SameSiteLaxMode,
	})

	// 로그인 성공 - 사용자 정보 + Access Token + Refresh Token 반환
	log.Printf("로그인 응답 생성: access_token=%s", accessToken[:20]+"...")
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success":       true,
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken, // Refresh Token도 응답에 포함
		"message":       "세션+JWT 기반 로그인 성공",
	}
	log.Printf("응답 데이터: %+v", response)
	json.NewEncoder(w).Encode(response)
}

// HandleLogout은 로그아웃 요청을 처리하는 핸들러 함수입니다
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// 세션 가져오기
	store := middleware.GetSessionStore()
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Printf("세션 가져오기 실패: %v", err)
		http.Error(w, "세션 오류", http.StatusInternalServerError)
		return
	}

	// 세션 무효화
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1 // 세션 즉시 삭제

	// 세션 저장
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

// HandleRefreshToken은 Refresh Token으로 Access Token을 갱신하는 핸들러 함수입니다
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("잘못된 HTTP 메서드: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var refreshToken string

	// 먼저 요청 본문에서 refresh_token 확인
	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err == nil {
		if token, ok := requestBody["refresh_token"].(string); ok && token != "" {
			refreshToken = token
		}
	}

	// 요청 본문에 없으면 쿠키에서 가져오기
	if refreshToken == "" {
		refreshTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			log.Printf("Refresh Token 없음: %v", err)
			http.Error(w, "Refresh token not found", http.StatusUnauthorized)
			return
		}
		refreshToken = refreshTokenCookie.Value
	}

	// Refresh Token으로 새로운 Access Token 생성
	newAccessToken, err := utills.RefreshAccessToken(refreshToken)
	if err != nil {
		log.Printf("Access Token 갱신 실패: %v", err)
		http.Error(w, "Token refresh failed", http.StatusUnauthorized)
		return
	}

	// Refresh Token에서 사용자 정보 추출
	userID, username, err := utills.GetUserInfoFromRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Refresh Token에서 사용자 정보 추출 실패: %v", err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// 새로운 Refresh Token 생성
	newRefreshToken, err := utills.GenerateRefreshToken(userID, username)
	if err != nil {
		log.Printf("새 Refresh Token 생성 실패: %v", err)
		// Access Token만 반환
		response := map[string]interface{}{
			"success":      true,
			"access_token": newAccessToken,
			"message":      "Access Token 갱신 성공",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 응답 생성
	response := map[string]interface{}{
		"success":       true,
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
		"message":       "Access Token 갱신 성공",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser 현재 로그인한 사용자 정보 조회
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// 세션에서 사용자 정보 가져오기
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
	var user struct {
		ID         int    `json:"id"`
		Userid     string `json:"userid"`
		Name       string `json:"name"`
		Role       string `json:"role"`
		Department string `json:"department"`
	}

	err = config.DB.QueryRow("SELECT id, userid, name, role, department FROM users WHERE id = ?", userID).Scan(
		&user.ID, &user.Userid, &user.Name, &user.Role, &user.Department)
	if err != nil {
		http.Error(w, "사용자 정보 조회 오류", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    user,
	})
}
