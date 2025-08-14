package model

import (
	"time"
)

// LeaveStatus 연차 현황 구조체
type LeaveStatus struct {
	UserID         int       `json:"user_id"`
	TotalEarned    int       `json:"total_earned_days"`
	UsedDays       int       `json:"used_days"`
	RemainingDays  int       `json:"remaining_days"`
	YearsOfService int       `json:"years_of_service"`
	HireDate       time.Time `json:"hire_date"`
}

// LeaveRequest 연차 신청 구조체
type LeaveRequest struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	LeaveType string    `json:"leave_type"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LeaveCalculation 연차 계산 결과 구조체
type LeaveCalculation struct {
	UserID          int       `json:"user_id"`
	HireDate        time.Time `json:"hire_date"`
	TargetDate      time.Time `json:"target_date"`
	YearsOfService  int       `json:"years_of_service"`
	TotalLeaveDays  int       `json:"total_leave_days"`
	CalculationType string    `json:"calculation_type"`
}

// LeavePolicy 연차 정책 구조체
type LeavePolicy struct {
	UnderOneYear map[string]interface{} `json:"under_one_year"`
	OverOneYear  map[string]interface{} `json:"over_one_year"`
}

// ==================== Response 구조체들 ====================

// LeaveStatusResponse 연차 현황 응답 구조체
type LeaveStatusResponse struct {
	Success bool         `json:"success"`
	Data    *LeaveStatus `json:"data"`
}

// LeaveCalculationResponse 연차 계산 응답 구조체
type LeaveCalculationResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

// LeaveUpdateResponse 연차 업데이트 응답 구조체
type LeaveUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Count   int    `json:"count"`
}
