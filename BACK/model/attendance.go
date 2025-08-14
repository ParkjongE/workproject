package model

import (
	"database/sql"
	"time"
)

// Attendance 출근 모델 (attendance_records 테이블에 맞춤)
type Attendance struct {
	RecordID     int            `json:"record_id" db:"record_id"`
	UserID       int            `json:"user_id" db:"user_id"`
	WorkDate     time.Time      `json:"work_date" db:"work_date"`
	CheckInTime  *time.Time     `json:"check_in_time" db:"check_in_time"`
	CheckOutTime *time.Time     `json:"check_out_time" db:"check_out_time"`
	Status       sql.NullString `json:"status" db:"status"` // NULL 허용
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
}

// AttendanceRequest 출근 요청 구조체 (기존)
type AttendanceRequest struct {
	UserID int    `json:"user_id"`
	Status string `json:"status"`
	Date   string `json:"date"`
}

// AttendanceResponse 출근 응답 구조체 (기존)
type AttendanceResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *Attendance `json:"data,omitempty"`
}

// 🔥 새로 추가할 구조체들

// CheckInRequest 출근 요청 구조체
type CheckInRequest struct {
	UserID      int    `json:"user_id"`
	Date        string `json:"date"`
	CheckInTime string `json:"check_in_time"`
}

// CheckOutRequest 퇴근 요청 구조체
type CheckOutRequest struct {
	UserID       int    `json:"user_id"`
	Date         string `json:"date"`
	CheckOutTime string `json:"check_out_time"`
}

// CheckInResponse 출근 응답 구조체
type CheckInResponse struct {
	Success     bool   `json:"success"`
	CheckInTime string `json:"check_in_time"`
	Message     string `json:"message"`
}

// CheckOutResponse 퇴근 응답 구조체
type CheckOutResponse struct {
	Success      bool   `json:"success"`
	CheckOutTime string `json:"check_out_time"`
	Message      string `json:"message"`
}

// AttendanceListResponse 출근 목록 응답 구조체 (통합)
type AttendanceListResponse struct {
	Date          string     `json:"date,omitempty"`           // "2006-01-02" 형식 (camelCase용)
	WorkDate      string     `json:"work_date,omitempty"`      // "2006-01-02" 형식 (snake_case용)
	CheckInTime   *time.Time `json:"checkInTime,omitempty"`    // 출근시간 (camelCase용)
	CheckInTime2  *time.Time `json:"check_in_time,omitempty"`  // 출근시간 (snake_case용)
	CheckOutTime  *time.Time `json:"checkOutTime,omitempty"`   // 퇴근시간 (camelCase용)
	CheckOutTime2 *time.Time `json:"check_out_time,omitempty"` // 퇴근시간 (snake_case용)
	WorkHours     float64    `json:"workHours,omitempty"`      // 근무시간 (camelCase용)
	WorkHours2    float64    `json:"work_hours,omitempty"`     // 근무시간 (snake_case용)
	Status        string     `json:"status"`                   // 상태
	WorkStatus    string     `json:"work_status,omitempty"`    // 근무 상태 (미완료, 지각, 정상, 초과근무)
}

// ToHistoryResponse Attendance를 AttendanceListResponse로 변환 (snake_case 형식)
func (a *Attendance) ToHistoryResponse() AttendanceListResponse {
	var status string
	if a.Status.Valid {
		status = a.Status.String
	}

	workHours := a.CalculateWorkHours()
	workStatus := a.GetWorkStatus()

	return AttendanceListResponse{
		WorkDate:      a.WorkDate.Format("2006-01-02"),
		CheckInTime2:  a.CheckInTime,
		CheckOutTime2: a.CheckOutTime,
		WorkHours2:    workHours,
		Status:        status,
		WorkStatus:    workStatus,
	}
}

// ToRecordsResponse Attendance를 AttendanceListResponse로 변환 (camelCase 형식)
func (a *Attendance) ToRecordsResponse() AttendanceListResponse {
	var status string
	if a.Status.Valid {
		status = a.Status.String
	}

	workHours := a.CalculateWorkHours()
	workStatus := a.GetWorkStatus()

	return AttendanceListResponse{
		Date:         a.WorkDate.Format("2006-01-02"),
		CheckInTime:  a.CheckInTime,
		CheckOutTime: a.CheckOutTime,
		WorkHours:    workHours,
		Status:       status,
		WorkStatus:   workStatus,
	}
}

// CalculateWorkHours 근무시간 계산 (시간 단위로 반환)
func (a *Attendance) CalculateWorkHours() float64 {
	if a.CheckInTime == nil || a.CheckOutTime == nil {
		return 0.0
	}

	duration := a.CheckOutTime.Sub(*a.CheckInTime)
	return duration.Hours()
}

// GetWorkStatus 근무 상태 반환 (미완료, 지각, 정상, 초과근무)
func (a *Attendance) GetWorkStatus() string {
	if a.CheckInTime == nil || a.CheckOutTime == nil {
		return "미완료"
	}

	workHours := a.CalculateWorkHours()
	if workHours >= 9 {
		return "초과근무"
	} else if workHours >= 8 {
		return "정상"
	} else {
		return "지각"
	}
}
