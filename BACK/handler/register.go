package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"BACK/model"
)

// RegisterService 사용자 등록 관련 서비스 인터페이스
type RegisterService interface {
	RegisterUser(user *model.User) error
	ValidateUserData(user *model.User) error
	CheckUserExists(userid string) error
	HashUserPassword(password string) (string, string, error)
}

// RegisterError 사용자 등록 관련 에러 구조체
type RegisterError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *RegisterError) Error() string {
	return e.Message
}

// RegisterHandler 사용자 등록 핸들러 구조체
type RegisterHandler struct {
	service RegisterService
}

// NewRegisterHandler 사용자 등록 핸들러 생성자
func NewRegisterHandler(service RegisterService) *RegisterHandler {
	return &RegisterHandler{service: service}
}

// ParseUserData 요청에서 사용자 데이터 파싱
func (h *RegisterHandler) ParseUserData(r *http.Request) (*model.User, error) {
	contentType := r.Header.Get("Content-Type")
	log.Printf("Content-Type: %s", contentType)

	var user model.User

	if strings.Contains(contentType, "multipart/form-data") {
		// multipart/form-data 처리
		log.Printf("multipart/form-data 형식으로 파싱")

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return nil, &RegisterError{Message: "폼 데이터 파싱 오류", Code: "FORM_PARSING_ERROR"}
		}

		user.Userid = r.FormValue("userid")
		user.Password = r.FormValue("password")
		user.Name = r.FormValue("name")
		user.Role = r.FormValue("role")
		user.Department = r.FormValue("department")

		// Form 데이터 처리 (기존 호환성)
		if hireDateStr := r.FormValue("hire_date"); hireDateStr != "" {
			if hireDate, err := time.Parse("2006-01-02", hireDateStr); err == nil {
				user.HireDate = hireDate
			}
		}

	} else if strings.Contains(contentType, "application/json") {
		// JSON 형식 처리
		log.Printf("JSON 형식으로 파싱")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, &RegisterError{Message: "요청 본문을 읽을 수 없습니다", Code: "BODY_READ_ERROR"}
		}

		if len(body) == 0 {
			return nil, &RegisterError{Message: "요청 본문이 비어 있습니다", Code: "EMPTY_BODY"}
		}

		// RegisterRequest로 파싱
		var reqData model.RegisterRequest
		if err := json.Unmarshal(body, &reqData); err != nil {
			log.Printf("JSON 파싱 오류: %v", err)
			return nil, &RegisterError{Message: fmt.Sprintf("잘못된 JSON 형식: %v", err), Code: "INVALID_JSON"}
		}

		// RegisterRequest를 User로 변환
		convertedUser, err := reqData.ToUser()
		if err != nil {
			log.Printf("날짜 파싱 오류: %v", err)
			return nil, &RegisterError{Message: fmt.Sprintf("입사날짜 형식이 올바르지 않습니다: %v", err), Code: "INVALID_DATE"}
		}

		user = *convertedUser
	} else {
		return nil, &RegisterError{Message: "지원하지 않는 Content-Type입니다. application/json 또는 multipart/form-data를 사용하세요.", Code: "UNSUPPORTED_CONTENT_TYPE"}
	}

	log.Printf("파싱된 사용자 데이터: userid='%s', password='%s', name='%s', role='%s', department='%s'",
		user.Userid, user.Password, user.Name, user.Role, user.Department)

	return &user, nil
}

// extractUserFieldsFromRawData rawData에서 사용자 필드 추출
func (h *RegisterHandler) extractUserFieldsFromRawData(user *model.User, rawData map[string]interface{}) {
	if userid, ok := rawData["userid"]; ok {
		if str, ok := userid.(string); ok {
			user.Userid = str
		}
	}
	if password, ok := rawData["password"]; ok {
		if str, ok := password.(string); ok {
			user.Password = str
		}
	}
	if name, ok := rawData["name"]; ok {
		if str, ok := name.(string); ok {
			user.Name = str
		}
	}
	if role, ok := rawData["role"]; ok {
		if str, ok := role.(string); ok {
			user.Role = str
		}
	}
	if department, ok := rawData["department"]; ok {
		if str, ok := department.(string); ok {
			user.Department = str
		}
	}
	if hireDate, ok := rawData["hire_date"]; ok {
		if str, ok := hireDate.(string); ok {
			if parsedDate, err := time.Parse("2006-01-02", str); err == nil {
				user.HireDate = parsedDate
			}
		}
	}
}

// HandleRegister 사용자 등록 요청을 처리하는 핸들러 함수
func (h *RegisterHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	log.Printf("회원가입 요청 메서드: %s, URL: %s", r.Method, r.URL.String())

	// GET 요청 처리
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

	// 사용자 데이터 파싱
	user, err := h.ParseUserData(r)
	if err != nil {
		if registerErr, ok := err.(*RegisterError); ok {
			switch registerErr.Code {
			case "FORM_PARSING_ERROR", "BODY_READ_ERROR", "EMPTY_BODY", "INVALID_JSON", "UNSUPPORTED_CONTENT_TYPE":
				http.Error(w, registerErr.Message, http.StatusBadRequest)
			default:
				http.Error(w, registerErr.Message, http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, "서버 오류", http.StatusInternalServerError)
		return
	}

	// 사용자 데이터 검증
	if err := h.service.ValidateUserData(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 사용자명 중복 확인
	if err := h.service.CheckUserExists(user.Userid); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// 비밀번호 해싱
	hashedPassword, salt, err := h.service.HashUserPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 해시된 비밀번호와 솔트 설정
	user.Password = hashedPassword
	user.Salt = salt

	// 사용자 등록
	if err := h.service.RegisterUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("새 사용자 등록 완료: %s (이름: %s, 역할: %s, 부서: %s)",
		user.Userid, user.Name, user.Role, user.Department)

	// 등록 성공 응답
	response := model.APIResponse{
		Success: true,
		Message: "회원가입이 완료되었습니다",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Register 사용자 등록 요청을 처리하는 핸들러 함수 (별칭)
func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.HandleRegister(w, r)
}
