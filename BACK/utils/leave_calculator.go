package utils

import (
	"time"
)

// LeaveStatus 사용자 연차 현황 구조체
type LeaveStatus struct {
	UserID        int `json:"user_id"`
	TotalEarned   int `json:"total_earned"`
	UsedDays      int `json:"used_days"`
	RemainingDays int `json:"remaining_days"`
}

// LeavePolicy 연차 정책 구조체
type LeavePolicy struct {
	MaxLeaveDays    int `json:"max_leave_days"`
	MinServiceYears int `json:"min_service_years"`
	LeaveIncrement  int `json:"leave_increment"`
}

// LeaveCalculator 연차 계산기 구조체
type LeaveCalculator struct{}

// NewLeaveCalculator 연차 계산기 생성자
func NewLeaveCalculator() *LeaveCalculator {
	return &LeaveCalculator{}
}

// CalculateUserLeaveStatus 사용자 연차 현황 계산
func (lc *LeaveCalculator) CalculateUserLeaveStatus(userID int, hireDate time.Time, usedDays int) *LeaveStatus {
	// 현재 연도
	currentYear := time.Now().Year()

	// 근무 연수 계산
	yearsOfService := currentYear - hireDate.Year()

	// 연차 계산 로직 (간단한 예시)
	totalEarned := 15 // 기본 15일
	if yearsOfService >= 1 {
		totalEarned = 15 + (yearsOfService * 2) // 1년 이상 근무 시 2일씩 추가
	}
	if totalEarned > 25 {
		totalEarned = 25 // 최대 25일
	}

	remainingDays := totalEarned - usedDays
	if remainingDays < 0 {
		remainingDays = 0
	}

	return &LeaveStatus{
		UserID:        userID,
		TotalEarned:   totalEarned,
		UsedDays:      usedDays,
		RemainingDays: remainingDays,
	}
}

// GetLeavePolicy 연차 정책 정보 조회
func (lc *LeaveCalculator) GetLeavePolicy() *LeavePolicy {
	return &LeavePolicy{
		MaxLeaveDays:    25,
		MinServiceYears: 1,
		LeaveIncrement:  2,
	}
}

// CalculateYearsOfService 근무 연수 계산
func (lc *LeaveCalculator) CalculateYearsOfService(hireDate, targetDate time.Time) int {
	return targetDate.Year() - hireDate.Year()
}

// CalculateTotalLeave 총 연차 계산
func (lc *LeaveCalculator) CalculateTotalLeave(hireDate, targetDate time.Time) int {
	yearsOfService := lc.CalculateYearsOfService(hireDate, targetDate)

	if yearsOfService < 1 {
		return 15 // 1년 미만 근무 시 15일
	}

	totalLeave := 15 + (yearsOfService * 2)
	if totalLeave > 25 {
		totalLeave = 25
	}

	return totalLeave
}
