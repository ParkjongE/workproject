package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"BACK/config"
	"BACK/models"
	"BACK/utills"
)

// HandleRegister는 사용자 등록 요청을 처리하는 핸들러 함수입니다
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	// 요청 메서드 로깅
	log.Printf("회원가입 요청 메서드: %s, URL: %s", r.Method, r.URL.String())

	// GET 요청 처리 - 회원가입 폼 요청으로 간주
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "회원가입 페이지입니다. POST 요청으로 회원가입을 진행하세요."})
		return
	}

	// POST 요청이 아닌 경우 (OPTIONS 제외)
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		log.Printf("잘못된 HTTP 메서드: %s (POST 요청이 필요합니다)", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// OPTIONS 요청 처리
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Content-Type 확인
	contentType := r.Header.Get("Content-Type")
	log.Printf("Content-Type: %s", contentType)

	// 사용자 객체 초기화
	var user models.User

	// Content-Type에 따라 다른 파싱 방법 사용
	if strings.Contains(contentType, "multipart/form-data") {
		// multipart/form-data 처리
		log.Printf("multipart/form-data 형식으로 파싱")

		// multipart 파서 생성 (최대 32MB)
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("multipart 파싱 오류: %v", err)
			http.Error(w, "폼 데이터 파싱 오류", http.StatusBadRequest)
			return
		}

		// 폼 데이터에서 필드 추출
		user.Userid = r.FormValue("userid")
		user.Password = r.FormValue("password")
		user.Name = r.FormValue("name")
		user.Role = r.FormValue("role")
		user.Department = r.FormValue("department")

		// hire_date 처리
		if hireDateStr := r.FormValue("hire_date"); hireDateStr != "" {
			if hireDate, err := time.Parse("2006-01-02", hireDateStr); err == nil {
				user.HireDate = hireDate
			}
		}

		log.Printf("multipart 파싱 결과: userid='%s', name='%s', role='%s', department='%s'",
			user.Userid, user.Name, user.Role, user.Department)

	} else if strings.Contains(contentType, "application/json") {
		// JSON 형식 처리
		log.Printf("JSON 형식으로 파싱")

		// 요청 본문 읽기
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("요청 본문 읽기 오류: %v", err)
			http.Error(w, "요청 본문을 읽을 수 없습니다", http.StatusBadRequest)
			return
		}

		log.Printf("원시 요청 본문 길이: %d", len(body))
		log.Printf("원시 요청 본문: %s", string(body))

		// 요청 본문이 비어있는지 확인
		if len(body) == 0 {
			log.Printf("요청 본문이 비어 있습니다")
			http.Error(w, "요청 본문이 비어 있습니다", http.StatusBadRequest)
			return
		}

		// JSON 파싱 시도
		if err := json.Unmarshal(body, &user); err != nil {
			log.Printf("JSON 파싱 오류: %v", err)
			log.Printf("파싱 실패한 JSON 데이터: %s", string(body))

			// 대체 파싱 방법 시도 - map으로 먼저 파싱
			var rawData map[string]interface{}
			if err := json.Unmarshal(body, &rawData); err != nil {
				log.Printf("대체 JSON 파싱도 실패: %v", err)
				log.Printf("rawData 파싱 실패한 JSON: %s", string(body))
				http.Error(w, fmt.Sprintf("잘못된 JSON 형식: %v", err), http.StatusBadRequest)
				return
			}

			log.Printf("rawData 파싱 성공: %+v", rawData)

			// 수동으로 필드 추출 (타입 안전하게)
			if userid, ok := rawData["userid"]; ok {
				if str, ok := userid.(string); ok {
					user.Userid = str
				} else {
					log.Printf("userid가 문자열이 아님: %T, 값: %v", userid, userid)
				}
			}
			if password, ok := rawData["password"]; ok {
				if str, ok := password.(string); ok {
					user.Password = str
				} else {
					log.Printf("password가 문자열이 아님: %T, 값: %v", password, password)
				}
			}
			if name, ok := rawData["name"]; ok {
				if str, ok := name.(string); ok {
					user.Name = str
				} else {
					log.Printf("name이 문자열이 아님: %T, 값: %v", name, name)
				}
			}
			if role, ok := rawData["role"]; ok {
				if str, ok := role.(string); ok {
					user.Role = str
				} else {
					log.Printf("role이 문자열이 아님: %T, 값: %v", role, role)
				}
			}
			if department, ok := rawData["department"]; ok {
				if str, ok := department.(string); ok {
					user.Department = str
				} else {
					log.Printf("department가 문자열이 아님: %T, 값: %v", department, department)
				}
			}
			if hireDate, ok := rawData["hire_date"]; ok {
				if str, ok := hireDate.(string); ok {
					if parsedDate, err := time.Parse("2006-01-02", str); err == nil {
						user.HireDate = parsedDate
					} else {
						log.Printf("hire_date 파싱 오류: %v", err)
					}
				} else {
					log.Printf("hire_date가 문자열이 아님: %T, 값: %v", hireDate, hireDate)
				}
			}
		}
	} else {
		log.Printf("지원하지 않는 Content-Type: %s", contentType)
		http.Error(w, "지원하지 않는 Content-Type입니다. application/json 또는 multipart/form-data를 사용하세요.", http.StatusBadRequest)
		return
	}

	// 디버깅을 위해 파싱된 사용자 데이터 출력
	log.Printf("파싱된 사용자 데이터: userid='%s', password='%s', name='%s', role='%s', department='%s'",
		user.Userid, user.Password, user.Name, user.Role, user.Department)

	// 필수 필드 검증
	if user.Userid == "" || user.Password == "" || user.Name == "" {
		log.Printf("필수 필드 누락: userid=%v, password=%v, name=%v",
			user.Userid != "", user.Password != "", user.Name != "")
		http.Error(w, "사용자명, 비밀번호, 이름은 필수 입력 항목입니다", http.StatusBadRequest)
		return
	}

	// hire_date가 비어있으면 기본값 설정 (2024-01-01)
	if user.HireDate.IsZero() {
		user.HireDate = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	// role이 비어있으면 기본값 "사원"으로 설정
	if user.Role == "" {
		user.Role = "사원"
	}

	// 사용자명 중복 확인
	var existingID int
	if err := config.DB.QueryRow("SELECT id FROM users WHERE userid = ?", user.Userid).Scan(&existingID); err == nil {
		// 사용자명이 이미 존재함
		http.Error(w, "이미 사용 중인 사용자명입니다", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		// 데이터베이스 오류
		log.Println("사용자 조회 오류:", err)
		http.Error(w, "서버 오류", http.StatusInternalServerError)
		return
	}

	// 1. 솔트 생성
	salt, err := utills.GenerateSalt(32)
	if err != nil {
		log.Println("솔트 생성 오류:", err)
		http.Error(w, "서버 오류", http.StatusInternalServerError)
		return
	}

	// 2. 비밀번호와 솔트를 결합하여 해싱
	hashedPassword, err := utills.HashPasswordWithSalt(user.Password, salt)
	if err != nil {
		log.Println("비밀번호 해싱 오류:", err)
		http.Error(w, "서버 오류", http.StatusInternalServerError)
		return
	}

	// 사용자 등록 - 데이터베이스에 저장
	log.Printf("데이터베이스에 사용자 저장 시도: %s", user.Userid)
	result, err := config.DB.Exec("INSERT INTO users (userid, password, salt, name, role, department, hire_date) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Userid, hashedPassword, salt, user.Name, user.Role, user.Department, user.HireDate)

	if err != nil {
		log.Printf("사용자 등록 오류: %v", err)
		http.Error(w, "사용자 등록에 실패했습니다", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("RowsAffected 확인 오류: %v", err)
	} else {
		log.Printf("영향받은 행 수: %d", rowsAffected)
	}

	log.Printf("새 사용자 등록 완료: %s (이름: %s, 역할: %s, 부서: %s)",
		user.Userid, user.Name, user.Role, user.Department)

	// 등록 성공 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
