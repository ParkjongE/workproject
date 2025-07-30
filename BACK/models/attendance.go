package models

import (
	"time"
)

// AttendanceRecord 출퇴근 기록 모델
type AttendanceRecord struct {
	RecordID     int        `json:"record_id"`
	UserID       int        `json:"user_id"`
	WorkDate     string     `json:"work_date"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	Status       *string    `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
}

// AttendanceRequest 출퇴근 요청 모델
type AttendanceRequest struct {
	UserID       int    `json:"user_id"`
	WorkDate     string `json:"work_date"`
	CheckInTime  string `json:"check_in_time,omitempty"`
	CheckOutTime string `json:"check_out_time,omitempty"`
}
