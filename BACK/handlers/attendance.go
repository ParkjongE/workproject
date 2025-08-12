package handlers

import (
	"BACK/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// AttendanceService 출근 관련 서비스 인터페이스
type AttendanceService interface {
	GetAttendance(userID int, startDate, endDate time.Time) ([]models.Attendance, error)
	PostAttendance(userID int, date time.Time, action string, timeStr string) error
}

// AttendanceHandler 출근 핸들러
type AttendanceHandler struct {
	service AttendanceService // 로컬 인터페이스 사용
}

// NewAttendanceHandler 출근 핸들러 생성자
func NewAttendanceHandler(service AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

// ==================== GET 메서드들 (GetAttendance 활용) ====================

// GetAttendanceHistory 출근 이력 조회 (최근 한 달)
func (h *AttendanceHandler) GetAttendanceHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("=== GetAttendanceHistory 시작 ===")

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 최근 30일 데이터 조회
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	log.Printf("조회 기간: %s ~ %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// GetAttendance 사용
	history, err := h.service.GetAttendance(userID, startDate, endDate)
	if err != nil {
		log.Printf("GetAttendance 호출 실패: %v", err)
		http.Error(w, "Failed to get attendance history", http.StatusInternalServerError)
		return
	}

	log.Printf("조회된 기록 수: %d", len(history))

	// 🔄 구조체 배열로 변환 (snake_case 형식)
	var records []models.AttendanceListResponse
	for _, attendance := range history {
		records = append(records, attendance.ToHistoryResponse())
	}

	h.sendJSONResponse(w, records)
	log.Println("=== GetAttendanceHistory 완료 ===")
}

// GetAttendanceByDate 특정 날짜의 출퇴근 상태 조회 (오늘 포함)
func (h *AttendanceHandler) GetAttendanceByDate(w http.ResponseWriter, r *http.Request) {
	log.Println("=== GetAttendanceByDate 시작 ===")
	log.Printf("요청 URL: %s", r.URL.String())

	userIDStr := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	log.Printf("쿼리 파라미터 - user_id: %s, date: %s", userIDStr, dateStr)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// GetAttendance 사용 (startDate와 endDate를 같게)
	history, err := h.service.GetAttendance(userID, date, date)
	if err != nil {
		log.Printf("GetAttendance 호출 실패: %v", err)
		// 해당 날짜에 기록이 없는 경우 빈 응답 (구조체 사용)
		emptyResponse := models.AttendanceListResponse{
			WorkDate:      dateStr,
			CheckInTime2:  nil,
			CheckOutTime2: nil,
			Status:        "",
		}
		h.sendJSONResponse(w, emptyResponse)
		return
	}

	// 결과가 있으면 첫 번째 항목 사용, 없으면 빈 응답
	var response models.AttendanceListResponse
	if len(history) > 0 {
		// 🎯 기존 ToHistoryResponse 메서드 활용 (snake_case 형식)
		response = history[0].ToHistoryResponse()
	} else {
		// 기록이 없는 경우 빈 응답
		response = models.AttendanceListResponse{
			WorkDate:      dateStr,
			CheckInTime2:  nil,
			CheckOutTime2: nil,
			Status:        "",
		}
	}

	log.Printf("GetAttendanceByDate 응답 - 날짜: %s, 출근시간: %v", response.WorkDate, response.CheckInTime2)
	h.sendJSONResponse(w, response)
	log.Println("=== GetAttendanceByDate 완료 ===")
}

// GetAttendanceRecords 사용자의 출퇴근 이력 조회 (URL 경로에서 user_id 추출)
func (h *AttendanceHandler) GetAttendanceRecords(w http.ResponseWriter, r *http.Request) {
	log.Println("=== GetAttendanceRecords 시작 ===")

	// URL에서 user_id 추출 (/api/attendance/records/{user_id})
	path := r.URL.Path
	log.Printf("요청 경로: %s", path)

	userIDStr := path[strings.LastIndex(path, "/")+1:]
	log.Printf("추출된 user_id: %s", userIDStr)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("user_id 변환 실패: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("변환된 user_id: %d", userID)

	// 최근 30일 데이터 조회
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	log.Printf("조회 기간: %s ~ %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// GetAttendance 사용
	history, err := h.service.GetAttendance(userID, startDate, endDate)
	if err != nil {
		log.Printf("GetAttendance 호출 실패: %v", err)
		http.Error(w, "Failed to get attendance history", http.StatusInternalServerError)
		return
	}

	log.Printf("조회된 기록 수: %d", len(history))

	// 🔄 구조체 배열로 변환 (camelCase 형식)
	var records []models.AttendanceListResponse
	for _, attendance := range history {
		records = append(records, attendance.ToRecordsResponse())
	}

	h.sendJSONResponse(w, records)
	log.Println("=== GetAttendanceRecords 완료 ===")
}

// ==================== POST 메서드들 (PostAttendance 활용) ====================

// CheckIn 출근 처리
func (h *AttendanceHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	log.Println("=== CheckIn 시작 ===")
	log.Printf("요청 메서드: %s, URL: %s", r.Method, r.URL.Path)

	// 🔄 models.CheckInRequest 사용
	var req models.CheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("요청 바디 파싱 실패: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("받은 요청 데이터: %+v", req)

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// 이미 해당 날짜에 출근 기록이 있는지 확인
	existing, err := h.service.GetAttendance(req.UserID, date, date)
	if err == nil && len(existing) > 0 && existing[0].CheckInTime != nil {
		log.Printf("이미 출근 기록이 있습니다: %v", existing[0].CheckInTime)
		http.Error(w, "이미 출근 기록이 있습니다", http.StatusConflict)
		return
	}

	// PostAttendance 사용
	err = h.service.PostAttendance(req.UserID, date, "checkin", req.CheckInTime)
	if err != nil {
		log.Printf("PostAttendance 호출 실패: %v", err)
		http.Error(w, "Failed to record check-in", http.StatusInternalServerError)
		return
	}

	// 🔄 models.CheckInResponse 사용
	response := models.CheckInResponse{
		Success:     true,
		CheckInTime: req.CheckInTime,
		Message:     "출근이 기록되었습니다.",
	}

	h.sendJSONResponse(w, response)
	log.Println("=== CheckIn 완료 ===")
}

// CheckOut 퇴근 처리
func (h *AttendanceHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	log.Println("=== CheckOut 시작 ===")

	// 🔄 models.CheckOutRequest 사용
	var req models.CheckOutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// PostAttendance 사용
	err = h.service.PostAttendance(req.UserID, date, "checkout", req.CheckOutTime)
	if err != nil {
		log.Printf("PostAttendance 호출 실패: %v", err)
		http.Error(w, "Failed to record check-out", http.StatusInternalServerError)
		return
	}

	// 🔄 models.CheckOutResponse 사용
	response := models.CheckOutResponse{
		Success:      true,
		CheckOutTime: req.CheckOutTime,
		Message:      "퇴근이 기록되었습니다.",
	}

	h.sendJSONResponse(w, response)
	log.Println("=== CheckOut 완료 ===")
}

// RecordAttendance 출근 기록 (기존 인터페이스 호환성)
func (h *AttendanceHandler) RecordAttendance(w http.ResponseWriter, r *http.Request) {
	log.Println("=== RecordAttendance 시작 ===")

	// 🔄 임시 구조체 대신 models.AttendanceRequest 사용
	var req models.AttendanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// PostAttendance 사용 (action을 status로 변환)
	var action string
	switch req.Status {
	case "출근":
		action = "checkin"
	case "퇴근":
		action = "checkout"
	default:
		action = "record"
	}

	err = h.service.PostAttendance(req.UserID, date, action, "")
	if err != nil {
		log.Printf("PostAttendance 호출 실패: %v", err)
		http.Error(w, "Failed to record attendance", http.StatusInternalServerError)
		return
	}

	// 🔄 models.AttendanceResponse 사용
	response := models.AttendanceResponse{
		Success: true,
		Message: "Attendance recorded successfully",
	}

	h.sendJSONResponse(w, response)
	log.Println("=== RecordAttendance 완료 ===")
}

// ==================== 공통 함수 ====================

// sendJSONResponse JSON 응답 전송 (공통 함수)
func (h *AttendanceHandler) sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON 인코딩 실패: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
